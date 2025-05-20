import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import Header from './components/Header';
import Leftbar from './components/Leftbar';
import Rightbar from './components/Rightbar';
import MatchRow from './components/MatchRow';

export default function App() {
  const [matches, setMatches] = useState([]);
  const [sport, setSport] = useState('soccer');
  const [loading, setLoading] = useState(false);
  const [mode, setMode] = useState('live');
  const [countdown, setCountdown] = useState(5);

  const fetchMatches = useCallback(async () => {
    setLoading(true);
    try {
      const endpoint = `/api/games?mode=${mode}&sport=${sport}`;
      const res = await axios.get(endpoint);
      const data = res.data ?? [];
      if (Array.isArray(data)) {
        setMatches(data);
      } else {
        console.warn('⚠️ API вернул не массив:', data);
        setMatches([]);
      }
    } catch (err) {
      console.error('Ошибка загрузки матчей:', err);
    } finally {
      setLoading(false);
      setCountdown(5);
    }
  }, [mode, sport]);

  useEffect(() => {
    fetchMatches();

    let refreshInterval;
    let countdownInterval;

    if (mode === 'live') {
      refreshInterval = setInterval(fetchMatches, 5000);
      countdownInterval = setInterval(() => {
        setCountdown((prev) => (prev > 0 ? prev - 1 : 5));
      }, 1000);
    }

    return () => {
      clearInterval(refreshInterval);
      clearInterval(countdownInterval);
    };
  }, [fetchMatches, mode]);

  return (
    <div className="min-h-screen bg-zinc-900 text-white flex flex-col">
      <Header mode={mode} setMode={setMode} onRefresh={fetchMatches} loading={loading} countdown={countdown} />
      <div className="flex flex-1 flex-col md:flex-row">
        <Leftbar selected={sport} onSelect={setSport} />
        <main className="flex-1 p-4 overflow-y-auto space-y-2">
          {loading && (
            <div className="text-center text-sm text-zinc-400 animate-pulse">Обновляется...</div>
          )}
          {matches.map((match) => (
            <MatchRow key={match.game_id} match={match} />
          ))}
        </main>
        <Rightbar />
      </div>
    </div>
  );
}
