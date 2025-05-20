package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Game struct {
	GameID     string `json:"game_id"`
	Time       string `json:"time"`
	TimeStatus string `json:"time_status"`
	League     string `json:"league"`
	Home       string `json:"home"`
	Away       string `json:"away"`
	Scores     string `json:"scores"`
	TimeStr    string `json:"time_str,omitempty"`
}

type Response struct {
	GamesLive []Game `json:"games_live"`
	GamesPre  []Game `json:"games_pre"`
}

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func fetchAndCacheGames(mode string, sport string) ([]Game, error) {
	login := os.Getenv("BOOKIES_LOGIN")
	token := os.Getenv("BOOKIES_TOKEN")
	task := "live"
	if mode == "pre" {
		task = "pre"
	}

	url := fmt.Sprintf(
		"https://bookiesapi.com/api/get.php?login=%s&token=%s&task=%s&bookmaker=bet365&sport=%s",
		login, token, task, sport,
	)

	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса API: %w", err)
	}

	body := resp.Body()
	if !json.Valid(body) {
		log.Println("⚠️ Ответ не является JSON:", string(body))
		return nil, fmt.Errorf("невалидный JSON от API: %s", string(body))
	}

	var parsed Response
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	var games []Game
	if mode == "pre" {
		games = parsed.GamesPre
	} else {
		games = parsed.GamesLive
	}

	for i := range games {
		sec, err := strconv.ParseInt(games[i].Time, 10, 64)
		if err == nil {
			games[i].TimeStr = time.Unix(sec, 0).Format("15:04")
		}
	}

	cacheKey := fmt.Sprintf("cached_games_%s_%s", mode, sport)
	data, _ := json.Marshal(games)
	err = rdb.Set(ctx, cacheKey, data, 30*time.Second).Err()
	if err != nil {
		log.Println("⚠️ Не удалось сохранить в Redis:", err)
	}

	log.Printf("🌐 [%s | %s] Загружены и закэшированы", mode, sport)
	return games, nil
}

func main() {
	godotenv.Load()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/games", func(c *gin.Context) {
		sport := c.DefaultQuery("sport", "soccer")
		mode := c.DefaultQuery("mode", "live")
		cacheKey := fmt.Sprintf("cached_games_%s_%s", mode, sport)

		cached, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedGames []Game
			if jsonErr := json.Unmarshal([]byte(cached), &cachedGames); jsonErr == nil {
				ttl, _ := rdb.TTL(ctx, cacheKey).Result()
				log.Printf("📦 [%s | %s] из Redis (TTL: %v сек)", mode, sport, int(ttl.Seconds()))
				c.JSON(http.StatusOK, cachedGames)
				return
			}
		}

		games, err := fetchAndCacheGames(mode, sport)
		if err != nil {
			log.Println("❌", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, games)
	})

	log.Println("🚀 Сервер запущен на http://localhost:8081")
	router.Run(":8081")
}
