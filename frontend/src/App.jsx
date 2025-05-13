
import { useEffect, useState } from "react";
import axios from "axios";
import MatchCard from "./MatchCard";
import Sidebar from "./Sidebar";

export default function App() {
  const [data, setData] = useState({});
  const [selectedSport, setSelectedSport] = useState("");
  const [selectedCategory, setSelectedCategory] = useState("Live");

  useEffect(() => {
    const fetchData = async () => {
      try {
        // const res = await axios.get("http://localhost:8081/structured");
        const res = await axios.get('/api/structured')
        setData(res.data);
        const firstSport = Object.keys(res.data)[0];
        setSelectedSport(firstSport);
      } catch (err) {
        console.error("Ошибка при получении данных:", err);
      }
    };

    fetchData();
  }, []);

  const sports = Object.keys(data);
  const currentMatches =
    data[selectedSport]?.[selectedCategory] || [];

  return (
    <div className="flex">
      <Sidebar
        sports={sports}
        selected={selectedSport}
        onSelect={setSelectedSport}
      />

      <main className="flex-1 p-4 bg-zinc-950 min-h-screen text-white">
        {selectedSport && (
          <>
            <div className="flex space-x-4 mb-6">
              {["Live", "Uncoming", "Popular"].map((cat) => (
                <button
                  key={cat}
                  onClick={() => setSelectedCategory(cat)}
                  className={`px-4 py-2 rounded-xl border ${
                    selectedCategory === cat
                      ? "bg-green-500 text-white"
                      : "bg-zinc-800 text-gray-300"
                  }`}
                >
                  {cat === "Uncoming" ? "Предстоящие" : cat}
                </button>
              ))}
            </div>

            {currentMatches.length === 0 ? (
              <p className="text-sm text-gray-500">Нет матчей</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {currentMatches.map((match, index) => (
                  <MatchCard
                    key={match.id || `${match.match || "?"}-${index}`}
                    match={match}
                  />
                ))}
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}
