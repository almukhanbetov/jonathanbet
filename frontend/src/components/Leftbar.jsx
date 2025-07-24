import React from 'react';

const sports = ['soccer', 'tennis', 'volleyball', 'basketball'];

export default function Leftbar({ selected, onSelect }) {
  return (
    <aside className="w-full md:w-48 bg-zinc-950 border-r border-zinc-800 p-4">
      <h2 className="font-semibold mb-4">Виды спорта</h2>
      <ul className="space-y-2">
        {sports.map((s) => (
          <li
            key={s}
            onClick={() => onSelect(s)}
            className={`cursor-pointer capitalize hover:text-green-400 ${
              selected === s ? 'text-green-400' : ''
            }`}
          >
            {s}
          </li>
        ))}
      </ul>
    </aside>
  );
}
