package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type RawItem struct {
	Type string `json:"type"`
}

type EV struct {
	Type string `json:"type"`
	ID   string `json:"ID"`
	NA   string `json:"NA"`
	SS   string `json:"SS"`
	TU   string `json:"TU"`
	TT   string `json:"TT"`
	CT   string `json:"CT"`
	CL   string `json:"CL"`
}

type PA struct {
	Type string `json:"type"`
	ID   string `json:"ID"`
	OD   string `json:"OD"`
	FI   string `json:"FI"`
	OR   string `json:"OR"`
}

type MatchCard struct {
	MatchID    string            `json:"id"`
	MatchName  string            `json:"match"`
	Time       string            `json:"time"`
	Score      string            `json:"score"`
	Tournament string            `json:"tournament"`
	Odds       map[string]string `json:"odds"`
}

type SportMatches struct {
	Popular  []MatchCard `json:"Popular"`
	Live     []MatchCard `json:"Live"`
	Uncoming []MatchCard `json:"Uncoming"`
}

var ctx = context.Background()

//	var rdb = redis.NewClient(&redis.Options{
//	    Addr: "localhost:6379",
//	})
var rdb = redis.NewClient(&redis.Options{
	Addr: "redis:6379", // <--- имя сервиса из docker-compose
})

func main() {
	godotenv.Load()
	startBackgroundUpdater()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/structured", func(c *gin.Context) {
		data, err := ParseStructuredMatches()
		if err != nil {
			log.Println("❌ Ошибка при получении данных:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	r.Run(":8081")
}

func startBackgroundUpdater() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("🔥 ПАНИКА в фоновом обновлении:", r)
			}
		}()

		for {
			err := fetchAndCacheData()
			if err != nil {
				log.Println("❌ [Фоновое обновление] Ошибка:", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

func fetchAndCacheData() error {
	client := resty.New()
	url := fmt.Sprintf("https://bookiesapi.com/api/get.php?login=%s&token=%s&task=bet365live",
		os.Getenv("BOOKIES_LOGIN"), os.Getenv("BOOKIES_TOKEN"))

	resp, err := client.R().Get(url)
	if err != nil {
		return fmt.Errorf("ошибка запроса к bookiesapi: %w", err)
	}

	if len(resp.Body()) < 50 {
		return fmt.Errorf("пустой или короткий ответ от API")
	}

	err = rdb.Set(ctx, "bet365live_json", resp.Body(), 35*time.Second).Err()
	if err != nil {
		return fmt.Errorf("не удалось сохранить в Redis: %w", err)
	}

	log.Println("📦 [Фоновое обновление] Данные успешно обновлены в Redis")
	return nil
}

func ParseStructuredMatches() (map[string]SportMatches, error) {
	cacheKey := "bet365live_json"
	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		log.Println("✅ Используем кэш из Redis")
		return parseFromJSON([]byte(cached))
	}
	log.Println("⚠️ Нет данных в Redis")
	return nil, fmt.Errorf("данные отсутствуют в кэше Redis")
}

func parseFromJSON(data []byte) (map[string]SportMatches, error) {
	var raw struct {
		Results [][]json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	clMap := make(map[string]string)
	matches := make([]EV, 0)
	odds := make([]PA, 0)

	for _, block := range raw.Results {
		for _, item := range block {
			var t RawItem
			json.Unmarshal(item, &t)

			switch t.Type {
			case "CL":
				var cl struct {
					Type string `json:"type"`
					ID   string `json:"ID"`
					NA   string `json:"NA"`
				}
				json.Unmarshal(item, &cl)
				clMap[cl.ID] = cl.NA

			case "EV":
				var ev EV
				json.Unmarshal(item, &ev)
				matches = append(matches, ev)

			case "PA":
				var pa PA
				json.Unmarshal(item, &pa)
				odds = append(odds, pa)
			}
		}
	}

	result := make(map[string]SportMatches)

	for _, match := range matches {
		if match.ID == "" || match.CL == "" {
			continue
		}

		sport := clMap[match.CL]
		if sport == "" {
			sport = "Other"
		}

		matchOdds := make(map[string]string)
		for _, o := range odds {
			if strings.HasPrefix(o.FI, match.ID) {
				switch o.OR {
				case "0":
					matchOdds["1"] = o.OD
				case "1":
					matchOdds["X"] = o.OD
				case "2":
					matchOdds["2"] = o.OD
				}
			}
		}

		card := MatchCard{
			MatchID:    match.ID,
			MatchName:  match.NA,
			Time:       parseTime(match.TU),
			Score:      match.SS,
			Tournament: match.CT,
			Odds:       matchOdds,
		}

		switch match.TT {
		case "1":
			result[sport] = appendTo(result[sport], "Live", card)
		case "0":
			result[sport] = appendTo(result[sport], "Uncoming", card)
		}
		result[sport] = appendTo(result[sport], "Popular", card)
	}

	for sport := range result {
		sort.Slice(result[sport].Popular, func(i, j int) bool {
			return result[sport].Popular[i].Time < result[sport].Popular[j].Time
		})
	}

	return result, nil
}

func appendTo(sm SportMatches, list string, match MatchCard) SportMatches {
	switch list {
	case "Live":
		sm.Live = append(sm.Live, match)
	case "Uncoming":
		sm.Uncoming = append(sm.Uncoming, match)
	case "Popular":
		sm.Popular = append(sm.Popular, match)
	}
	return sm
}

func parseTime(raw string) string {
	if len(raw) < 12 {
		return raw
	}
	t, err := time.Parse("20060102150405", raw)
	if err != nil {
		return raw
	}
	return t.Format("2006-01-02 15:04")
}
