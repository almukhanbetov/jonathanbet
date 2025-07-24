import React from 'react';

export default function MatchRow({ match }) {
  return (
    <div className="bg-zinc-800 p-3 rounded flex flex-col md:flex-row justify-between items-start md:items-center shadow hover:bg-zinc-700 transition">
      <div>
        <div className="text-sm text-zinc-400">{match.league}</div>
        <div className="font-semibold">
          {match.home} vs {match.away}
        </div>
      </div>
      <div className="text-right mt-2 md:mt-0">
        <div className="text-green-400 font-bold">{match.scores}</div>
        <div className="text-xs text-zinc-400">{match.time_str}</div>
      </div>
    </div>
  );
}
