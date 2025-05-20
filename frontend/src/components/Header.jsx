import React from 'react';

export default function Header({ mode, setMode, onRefresh, loading, countdown }) {
  return (
    <header className="bg-zinc-800 px-6 py-4 shadow flex items-center justify-between">
      <div className="flex items-center gap-4">
        <h1 className="text-xl font-bold">JonathanBet</h1>
        <button
          onClick={() => setMode('live')}
          className={`px-4 py-2 rounded ${
            mode === 'live' ? 'bg-green-600' : 'bg-zinc-700'
          }`}
        >
          Live
        </button>
        <button
          onClick={() => setMode('pre')}
          className={`px-4 py-2 rounded ${
            mode === 'pre' ? 'bg-green-600' : 'bg-zinc-700'
          }`}
        >
          Предстоящие
        </button>
        <button
          onClick={onRefresh}
          className="px-4 py-2 rounded bg-zinc-700 hover:bg-zinc-600"
        >
          {loading ? '🔄 Обновляется...' : `🔁 Обновить (${countdown})`}
        </button>
      </div>
      <div className="flex gap-4">
        <button className="bg-zinc-700 px-4 py-2 rounded">Войти</button>
        <button className="bg-green-600 px-4 py-2 rounded">Корзина</button>
      </div>
    </header>
  );
}
