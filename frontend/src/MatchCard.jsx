
export default function MatchCard({ match }) {
  return (
    <div className="bg-zinc-900 p-4 rounded-2xl shadow-md">
      <h2 className="text-lg font-semibold">{match.match}</h2>
      <p className="text-sm text-gray-400">Начало: {match.time}</p>
      <p className="text-sm text-gray-400">Счёт: {match.score}</p>
      <p className="text-sm text-gray-400 mb-2">Турнир: {match.tournament}</p>

      <div className="grid grid-cols-3 gap-2 mt-3">
        {match.odds &&
          Object.entries(match.odds).map(([label, odd]) => (
            <div
              key={label}
              className="bg-zinc-800 p-2 rounded-lg text-center border border-zinc-700"
            >
              <div className="text-sm">{label}</div>
              <div className="font-bold text-green-400">{odd}</div>
            </div>
          ))}
      </div>
    </div>
  );
}
