export default function Sidebar({ sports, selected, onSelect }) {
  return (
    <aside className="w-64 bg-zinc-950 text-white p-4 border-r border-zinc-800">
      <h2 className="text-xl font-bold mb-4">Спорт</h2>
      <ul className="space-y-2">
        {sports.map((sport) => (
          <li
            key={sport}
            className={`cursor-pointer hover:text-green-400 ${sport === selected ? "text-green-400 font-bold" : ""}`}
            onClick={() => onSelect(sport)}
          >
            {sport}
          </li>
        ))}
      </ul>
    </aside>
  );
}
