import { useState, useEffect, useCallback } from 'react';
import { Folder, FolderOpen, ChevronRight, Home, HardDrive, CornerUpLeft, X } from 'lucide-react';
import { api, FSEntry } from '../api-client';

interface FileSystemBrowserProps {
  selectFolder?: boolean;
  onSelect?: (path: string) => void;
  onClose?: () => void;
}

export default function FileSystemBrowser({
  selectFolder = true,
  onSelect,
  onClose
}: FileSystemBrowserProps) {
  const [currentPath, setCurrentPath] = useState('');
  const [entries, setEntries] = useState<FSEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [homeDir, setHomeDir] = useState('');
  const [pathHistory, setPathHistory] = useState<string[]>([]);

  const loadDirectory = useCallback(async (path: string, addToHistory = true) => {
    setLoading(true);
    setError(null);
    try {
      const data = await api.browseFS(path);
      setEntries(data.entries);
      setCurrentPath(data.cwd);
      if (addToHistory && pathHistory[pathHistory.length - 1] !== data.cwd) {
        setPathHistory(prev => [...prev, data.cwd]);
      }
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, [pathHistory]);

  const goBack = useCallback(() => {
    if (pathHistory.length > 1) {
      const newHistory = pathHistory.slice(0, -1);
      const prevPath = newHistory[newHistory.length - 1];
      setPathHistory(newHistory);
      loadDirectory(prevPath, false);
    }
  }, [pathHistory, loadDirectory]);

  const loadHomeDir = useCallback(async () => {
    setPathHistory([]);
    try {
      const data = await api.getHomeDir();
      setHomeDir(data.home);
      await loadDirectory(data.home);
    } catch (e: any) {
      console.error('Failed to load home dir:', e);
      setError('Failed to load home directory');
    }
  }, [loadDirectory]);

  const navigateTo = (path: string) => {
    loadDirectory(path);
  };

  const handleEntryClick = (entry: FSEntry) => {
    if (entry.is_dir) {
      navigateTo(entry.path);
    }
  };

  const handleEntryDoubleClick = (entry: FSEntry) => {
    if (!entry.is_dir) {
      onSelect?.(entry.path);
      onClose?.();
    }
  };

  const handleSelectCurrent = () => {
    onSelect?.(currentPath);
    onClose?.();
  };

  const closeModal = () => {
    onClose?.();
  };

  const formatSize = (size?: number): string => {
    if (!size) return '';
    if (size < 1024) return `${size} B`;
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
    if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(1)} MB`;
    return `${(size / (1024 * 1024 * 1024)).toFixed(1)} GB`;
  };

  // Handle keyboard escape
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') closeModal();
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, []);

  useEffect(() => {
    loadHomeDir();
  }, [loadHomeDir]);

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      style={{ backgroundColor: 'rgba(0, 0, 0, 0.5)' }}
      onClick={() => closeModal()}
    >
      <div
        className="card w-full max-w-2xl max-h-[80vh] flex flex-col"
        style={{ backgroundColor: 'var(--bg-card)' }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b" style={{ borderColor: 'var(--border-color)' }}>
          <div className="flex items-center gap-2">
            <FolderOpen size={20} style={{ color: 'var(--accent)' }} />
            <h3 className="font-semibold" style={{ color: 'var(--text-primary)' }}>Browse Folder</h3>
          </div>
          <button className="btn-icon w-8 h-8" onClick={() => closeModal()} title="Close (Esc)">
            <X size={20} />
          </button>
        </div>

        {/* Navigation Bar */}
        <div className="flex items-center gap-2 p-4 border-b" style={{ borderColor: 'var(--border-color)' }}>
          <button className="btn-icon" onClick={() => goBack()} disabled={pathHistory.length <= 1} title="Back">
            <CornerUpLeft size={16} style={{ transform: 'rotate(180deg)' }} />
          </button>
          <button className="btn-icon" onClick={() => loadHomeDir()} title="Home">
            <Home size={16} />
          </button>
          {homeDir && (
            <button className="btn-icon" onClick={() => navigateTo(homeDir)} title="Home Directory">
              <HardDrive size={16} />
            </button>
          )}
          <div
            className="flex-1 flex items-center gap-1 px-3 py-2 rounded-lg font-mono text-sm"
            style={{ backgroundColor: 'var(--bg-muted)', color: 'var(--text-secondary)' }}
          >
            <span className="truncate">{currentPath || 'Loading...'}</span>
          </div>
        </div>

        {/* Content */}
        {error ? (
          <div className="p-4 text-center" style={{ color: 'var(--destructive)' }}>
            <p>{error}</p>
            <button className="btn-secondary mt-2" onClick={loadHomeDir}>Go to Home</button>
          </div>
        ) : loading ? (
          <div className="p-8 text-center" style={{ color: 'var(--text-muted)' }}>
            <p>Loading...</p>
          </div>
        ) : (
          <div className="flex-1 overflow-y-auto p-2">
            {entries.length === 0 ? (
              <p className="text-center py-4" style={{ color: 'var(--text-muted)' }}>No items in this folder</p>
            ) : (
              <div className="space-y-1">
                {entries.map((entry) => (
                  <button
                    key={entry.path}
                    className="w-full flex items-center gap-3 p-3 rounded-lg hover:border-[var(--accent)] transition-all text-left"
                    style={{
                      backgroundColor: entry.is_dir ? 'var(--bg-muted)' : 'var(--bg-card)',
                      border: '1px solid var(--border-subtle)'
                    }}
                    onClick={() => handleEntryClick(entry)}
                    onDoubleClick={() => handleEntryDoubleClick(entry)}
                  >
                    {entry.is_dir ? (
                      <>
                        <Folder size={18} style={{ color: 'var(--accent)' }} />
                        <span className="font-medium" style={{ color: 'var(--text-primary)' }}>{entry.name}</span>
                        <ChevronRight size={16} style={{ color: 'var(--text-muted)', marginLeft: 'auto' }} />
                      </>
                    ) : (
                      <>
                        <ChevronRight size={18} style={{ color: 'var(--text-muted)' }} />
                        <span className="font-medium" style={{ color: 'var(--text-primary)' }}>{entry.name}</span>
                        {entry.size && (
                          <span className="text-xs" style={{ color: 'var(--text-muted)', marginLeft: 'auto' }}>
                            {formatSize(entry.size)}
                          </span>
                        )}
                      </>
                    )}
                  </button>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Footer Actions */}
        {selectFolder && (
          <div className="flex items-center justify-end gap-2 p-4 border-t" style={{ borderColor: 'var(--border-color)' }}>
            <button className="btn-secondary" onClick={() => closeModal()}>Cancel</button>
            <button className="btn-primary" onClick={() => handleSelectCurrent()}>
              Select This Folder
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
