import React, { useState, useEffect } from 'react';
import './App.css';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const CYCLE_DURATION_MS = parseInt(import.meta.env.VITE_CYCLE_DURATION_MS || '18000000', 10);

// WindowPlayer component handles the logic for a single display window
const WindowPlayer = ({ windowData, playlist, mediaList, globalSync }) => {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [cycleBlock, setCycleBlock] = useState(() => Math.floor(Date.now() / CYCLE_DURATION_MS));

  // Determine what media to show
  let activeMedia = null;
  let isSyncing = false;

  if (globalSync && globalSync.active && globalSync.media_id) {
    activeMedia = mediaList.find(m => m.id === globalSync.media_id);
    isSyncing = true;
  } else if (playlist && playlist.media_ids && playlist.media_ids.length > 0) {
    const safeIndex = currentIndex % playlist.media_ids.length;
    const mediaId = playlist.media_ids[safeIndex];
    activeMedia = mediaList.find(m => m.id === mediaId);
  }

  // 5-Hour Cycle Master Clock
  useEffect(() => {
    const interval = setInterval(() => {
      const currentBlock = Math.floor(Date.now() / CYCLE_DURATION_MS);
      if (currentBlock !== cycleBlock) {
        setCycleBlock(currentBlock);
        setCurrentIndex(0); // Reset index on boundary
      }
    }, 1000);
    return () => clearInterval(interval);
  }, [cycleBlock]);

  // Handle media progression for images/fallback/blank
  useEffect(() => {
    if (isSyncing || !activeMedia) return;

    if (activeMedia.type === 'video') return;

    const durationMs = (activeMedia.duration || 5) * 1000;
    const timer = setTimeout(() => {
      setCurrentIndex(prev => prev + 1);
    }, durationMs);

    return () => clearTimeout(timer); // cleanup prevents stale timeouts
  }, [activeMedia, isSyncing, currentIndex, cycleBlock]);

  const handleVideoEnded = () => {
    if (isSyncing) return;
    setCurrentIndex(prev => prev + 1);
  };

  const handleError = (e) => {
    e.target.style.display = 'none';
    const fallback = e.target.parentElement.querySelector('.fallback-media');
    if (fallback) fallback.style.display = 'block';
    
    if (!isSyncing && activeMedia && activeMedia.type === 'video') {
       setTimeout(handleVideoEnded, (activeMedia.duration || 5) * 1000);
    }
  };

  return (
    <div className="window-panel">
      <div className="window-header">
        <h2>{windowData.name}</h2>
        <p>Status: {isSyncing ? 'SYNC OVERRIDE' : windowData.status}</p>
        <p style={{fontSize: '0.8em', color: 'var(--text-muted)'}}>Cycle Block: {cycleBlock}</p>
      </div>
      
      <div className="media-container">
        {isSyncing && <div className="sync-indicator">SYNC ACTIVE</div>}
        
        {!activeMedia ? (
          <div className="fallback-media">
            <h3>No Media</h3>
            <p>Playlist empty or media not found.</p>
          </div>
        ) : activeMedia.type === 'blank' ? (
           <div 
             key={`${activeMedia.id}-${cycleBlock}-${isSyncing ? 'sync' : 'local'}`} 
             style={{ width: '100%', height: '100%', background: 'black' }} 
           />
        ) : activeMedia.type === 'video' ? (
          <video 
            key={`${activeMedia.id}-${cycleBlock}-${isSyncing ? 'sync' : 'local'}`} // force remount on cycle cross or sync
            src={`/media/${activeMedia.file_path.split('/').pop()}`}
            className="media-content"
            autoPlay
            muted
            onEnded={handleVideoEnded}
            onError={handleError}
          />
        ) : (
          <img 
            key={`${activeMedia.id}-${cycleBlock}-${isSyncing ? 'sync' : 'local'}`}
            src={`/media/${activeMedia.file_path.split('/').pop()}`}
            className="media-content"
            alt={activeMedia.name}
            onError={handleError}
          />
        )}
        
        {/* Fallback Display */}
        <div className="fallback-media" style={{ display: 'none', position: 'absolute' }}>
          <h3>{activeMedia?.name || 'Unknown'}</h3>
          <p>[{activeMedia?.type?.toUpperCase()}] ID: {activeMedia?.id}</p>
          <p>Duration: {activeMedia?.duration}s</p>
        </div>
      </div>
    </div>
  );
};

function App() {
  const [windows, setWindows] = useState([]);
  const [playlists, setPlaylists] = useState([]);
  const [mediaList, setMediaList] = useState([]);
  const [syncState, setSyncState] = useState({ active: false });
  
  const [syncInputId, setSyncInputId] = useState('M2');
  const [syncInputDuration, setSyncInputDuration] = useState('10');

  const [addMediaWindowId, setAddMediaWindowId] = useState('');
  const [addMediaId, setAddMediaId] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [winRes, playRes, medRes] = await Promise.all([
          fetch(`${API_BASE}/windows`),
          fetch(`${API_BASE}/playlists`),
          fetch(`${API_BASE}/media`)
        ]);
        
        const fetchedWindows = await winRes.json() || [];
        setWindows(fetchedWindows);
        setPlaylists(await playRes.json() || []);
        setMediaList(await medRes.json() || []);

        if (fetchedWindows.length > 0) {
          setAddMediaWindowId(fetchedWindows[0].id);
        }
      } catch (err) {
        console.error("Failed to fetch initial data", err);
      }
    };
    
    fetchData();

    // Setup polling for sync state every 1 second
    const interval = setInterval(async () => {
      try {
        const res = await fetch(`${API_BASE}/sync`);
        if (res.ok) {
          const data = await res.json();
          setSyncState(data);
        }
      } catch (err) {}
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  const handleTriggerSync = async (e) => {
    e.preventDefault();
    try {
      await fetch(`${API_BASE}/sync`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          media_id: syncInputId,
          duration: parseInt(syncInputDuration, 10)
        })
      });
    } catch (err) {
      alert("Failed to trigger sync");
    }
  };

  const handleAddMedia = async (e) => {
    e.preventDefault();
    const win = windows.find(w => w.id === addMediaWindowId);
    if (!win) return;

    try {
      const res = await fetch(`${API_BASE}/playlists/${win.playlist_id}/media`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ media_id: addMediaId })
      });

      if (!res.ok) {
        const errData = await res.json();
        alert(`Error: ${errData.error}`);
        return;
      }

      const updatedPlaylist = await res.json();
      
      // Update state without refreshing
      setPlaylists(prev => prev.map(p => p.id === updatedPlaylist.id ? updatedPlaylist : p));
      setAddMediaId(''); // clear input
    } catch (err) {
      alert("Failed to add media");
    }
  };

  return (
    <div className="app-container">
      <div className="header" style={{ flexDirection: 'column', alignItems: 'flex-start', gap: '1rem' }}>
        <h1>Media Sync Player</h1>
        
        <div style={{ display: 'flex', gap: '2rem', width: '100%', flexWrap: 'wrap' }}>
          {/* Sync Form */}
          <form className="sync-controls" onSubmit={handleTriggerSync}>
            <label>Sync Media ID:</label>
            <input 
              type="text" 
              value={syncInputId}
              onChange={e => setSyncInputId(e.target.value)}
              required 
            />
            <label>Duration (s):</label>
            <input 
              type="number" 
              value={syncInputDuration}
              onChange={e => setSyncInputDuration(e.target.value)}
              min="1"
              required 
            />
            <button type="submit">Trigger Global Sync</button>
          </form>

          {/* Add Media Form */}
          <form className="sync-controls" onSubmit={handleAddMedia}>
            <label>Add Media to Window:</label>
            <select style={{padding: '0.5rem'}} value={addMediaWindowId} onChange={e => setAddMediaWindowId(e.target.value)}>
              {windows.map(w => (
                <option key={w.id} value={w.id}>{w.name}</option>
              ))}
            </select>
            <label>Media ID:</label>
            <input 
              type="text" 
              value={addMediaId}
              onChange={e => setAddMediaId(e.target.value)}
              required 
            />
            <button type="submit">Add Media</button>
          </form>
        </div>
      </div>

      <div className="windows-grid">
        {windows.map(w => (
          <WindowPlayer 
            key={w.id}
            windowData={w}
            playlist={playlists.find(p => p.id === w.playlist_id)}
            mediaList={mediaList}
            globalSync={syncState}
          />
        ))}
      </div>
    </div>
  );
}

export default App;
