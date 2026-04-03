import { useState, useEffect, useCallback } from 'react';
import { Folder, FolderOpen, ChevronRight, Home, HardDrive, CornerUpLeft } from 'lucide-react';
import { api, FSEntry } from '../api-client';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog';

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
      if (addToHistory) {
        setPathHistory(prev => {
          // Only add to history if different from last entry
          if (prev.length > 0 && prev[prev.length - 1] === data.cwd) {
            return prev;
          }
          return [...prev, data.cwd];
        });
      }
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  const goBack = useCallback(() => {
    setPathHistory(prev => {
      if (prev.length > 1) {
        const newHistory = prev.slice(0, -1);
        const prevPath = newHistory[newHistory.length - 1];
        // Load previous directory without adding to history
        (async () => {
          setLoading(true);
          setError(null);
          try {
            const data = await api.browseFS(prevPath);
            setEntries(data.entries);
            setCurrentPath(data.cwd);
          } catch (e: any) {
            setError(e.message);
          } finally {
            setLoading(false);
          }
        })();
        return newHistory;
      }
      return prev;
    });
  }, []);

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

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      onClose?.();
    }
  };

  const formatSize = (size?: number): string => {
    if (!size) return '';
    if (size < 1024) return `${size} B`;
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
    if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(1)} MB`;
    return `${(size / (1024 * 1024 * 1024)).toFixed(1)} GB`;
  };

  useEffect(() => {
    loadHomeDir();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <Dialog defaultOpen onOpenChange={handleOpenChange}>
      <DialogContent className="max-w-2xl max-h-[80vh] flex flex-col">
        {/* Header */}
        <DialogHeader>
          <div className="flex items-center gap-2">
            <FolderOpen size={20} className="text-primary" />
            <DialogTitle>Browse Folder</DialogTitle>
          </div>
        </DialogHeader>

        {/* Navigation Bar */}
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="icon" onClick={() => goBack()} disabled={pathHistory.length <= 1} title="Back">
            <CornerUpLeft size={16} className="rotate-180" />
          </Button>
          <Button variant="ghost" size="icon" onClick={() => loadHomeDir()} title="Home">
            <Home size={16} />
          </Button>
          {homeDir && (
            <Button variant="ghost" size="icon" onClick={() => navigateTo(homeDir)} title="Home Directory">
              <HardDrive size={16} />
            </Button>
          )}
          <div
            className="flex-1 flex items-center gap-1 px-3 py-2 rounded-md font-mono text-sm bg-muted text-muted-foreground"
          >
            <span className="truncate">{currentPath || 'Loading...'}</span>
          </div>
        </div>

        {/* Content */}
        {error ? (
          <div className="p-4 text-center text-destructive">
            <p>{error}</p>
            <Button variant="secondary" className="mt-2" onClick={loadHomeDir}>Go to Home</Button>
          </div>
        ) : loading ? (
          <div className="p-8 text-center text-muted-foreground">
            <p>Loading...</p>
          </div>
        ) : (
          <div className="flex-1 overflow-y-auto p-2">
            {entries.length === 0 ? (
              <p className="text-center py-4 text-muted-foreground">No items in this folder</p>
            ) : (
              <div className="space-y-1">
                {entries.map((entry) => (
                  <button
                    key={entry.path}
                    className="w-full flex items-center gap-3 p-3 rounded-md hover:border-primary transition-all text-left bg-card border border-border/50"
                    onClick={() => handleEntryClick(entry)}
                    onDoubleClick={() => handleEntryDoubleClick(entry)}
                  >
                    {entry.is_dir ? (
                      <>
                        <Folder size={18} className="text-primary" />
                        <span className="font-medium text-foreground">{entry.name}</span>
                        <ChevronRight size={16} className="text-muted-foreground ml-auto" />
                      </>
                    ) : (
                      <>
                        <ChevronRight size={18} className="text-muted-foreground" />
                        <span className="font-medium text-foreground">{entry.name}</span>
                        {entry.size && (
                          <span className="text-xs text-muted-foreground ml-auto">
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
          <DialogFooter>
            <Button variant="outline" onClick={() => onClose?.()}>Cancel</Button>
            <Button onClick={handleSelectCurrent}>
              Select This Folder
            </Button>
          </DialogFooter>
        )}
      </DialogContent>
    </Dialog>
  );
}
