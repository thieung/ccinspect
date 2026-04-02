import { useState, useEffect } from 'react';
import { api, ConfigData } from '../api-client';
import { Settings, Save, RefreshCw, CheckCircle2, FolderOpen, Plus, Trash2, X } from 'lucide-react';
import FileSystemBrowser from '../components/FileSystemBrowser';

export default function SettingsView() {
  const [config, setConfig] = useState<ConfigData | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  // Form fields
  const [scanPaths, setScanPaths] = useState<string[]>([]);
  const [excludePaths, setExcludePaths] = useState('');
  const [maxDepth, setMaxDepth] = useState(5);
  const [defaultOutput, setDefaultOutput] = useState('table');

  // File browser
  const [showBrowser, setShowBrowser] = useState(false);
  const [browserTargetPath, setBrowserTargetPath] = useState<'scan' | 'exclude'>('scan');

  async function loadConfig() {
    setLoading(true);
    setError(null);
    try {
      const data = await api.getConfig();
      setConfig(data.config);
      setScanPaths(data.config?.scan_paths ?? []);
      setExcludePaths(data.config?.exclude_paths?.join('\n') ?? '');
      setMaxDepth(data.config?.max_depth ?? 5);
      setDefaultOutput(data.config?.default_output ?? 'table');
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  async function saveConfig() {
    setSaving(true);
    setError(null);
    setSuccess(false);
    try {
      await api.saveConfig({
        scan_paths: scanPaths,
        exclude_paths: excludePaths.split('\n').map(s => s.trim()).filter(Boolean),
        max_depth: maxDepth,
        default_output: defaultOutput,
      });
      setSuccess(true);
      setTimeout(() => setSuccess(false), 3000);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setSaving(false);
    }
  }

  function openBrowser(target: 'scan' | 'exclude' = 'scan') {
    setBrowserTargetPath(target);
    setShowBrowser(true);
  }

  function handlePathSelect(path: string) {
    if (browserTargetPath === 'scan') {
      if (!scanPaths.includes(path)) {
        setScanPaths([...scanPaths, path]);
      }
    } else if (browserTargetPath === 'exclude') {
      setExcludePaths(excludePaths ? `${excludePaths}\n${path}` : path);
    }
    setShowBrowser(false);
  }

  function removeScanPath(pathToRemove: string) {
    setScanPaths(scanPaths.filter(p => p !== pathToRemove));
  }

  function removeExcludePath(index: number) {
    const lines = excludePaths.split('\n');
    lines.splice(index, 1);
    setExcludePaths(lines.join('\n'));
  }

  useEffect(() => {
    loadConfig();
  }, []);

  const excludePathLines = excludePaths.split('\n').filter(Boolean);

  return (
    <div className="space-y-6">
      <div className="card p-6">
        <div className="flex items-center gap-3 mb-6">
          <div className="w-12 h-12 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'var(--accent)' }}>
            <Settings size={24} className="text-white" />
          </div>
          <div>
            <h2 className="text-xl font-semibold" style={{ color: 'var(--text-primary)' }}>Settings</h2>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Configure ccinspect scanning behavior</p>
          </div>
        </div>

        {loading ? (
          <div className="text-center py-12" style={{ color: 'var(--text-muted)' }}>
            <RefreshCw size={32} className="animate-spin mx-auto mb-4" />
            <p>Loading config...</p>
          </div>
        ) : error && !config ? (
          <div className="card p-4 border-[var(--destructive)]">
            <p className="text-sm" style={{ color: 'var(--destructive)' }}>Error: {error}</p>
            <button onClick={loadConfig} className="btn-secondary mt-2 text-sm">Retry</button>
          </div>
        ) : (
          <div className="space-y-6">
            {/* Scan Paths */}
            <div>
              <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-primary)' }}>Scan Paths</label>
              <div className="space-y-2">
                {scanPaths.map((path) => (
                  <div key={path} className="flex items-center gap-2">
                    <div
                      className="flex-1 flex items-center gap-2 px-3 py-2 rounded-lg font-mono text-sm"
                      style={{ backgroundColor: 'var(--bg-muted)', color: 'var(--text-secondary)', border: '1px solid var(--border-subtle)' }}
                    >
                      <FolderOpen size={14} style={{ color: 'var(--accent)' }} />
                      <span className="truncate">{path}</span>
                    </div>
                    <button className="btn-icon" onClick={() => removeScanPath(path)} title="Remove">
                      <X size={16} />
                    </button>
                  </div>
                ))}
                <button onClick={() => openBrowser('scan')} className="btn-secondary flex items-center gap-2 w-full">
                  <Plus size={16} />
                  Add Folder...
                </button>
              </div>
              <p className="text-xs mt-1" style={{ color: 'var(--text-muted)' }}>Select directories to scan for .claude/ directories.</p>
            </div>

            {/* Exclude Paths */}
            <div>
              <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-primary)' }}>Exclude Paths</label>
              <div className="space-y-2">
                {excludePathLines.map((line, index) => (
                  <div key={`${line}-${index}`} className="flex items-center gap-2">
                    <div
                      className="flex-1 flex items-center gap-2 px-3 py-2 rounded-lg font-mono text-sm"
                      style={{ backgroundColor: 'var(--bg-muted)', color: 'var(--text-secondary)', border: '1px solid var(--border-subtle)' }}
                    >
                      <Trash2 size={14} style={{ color: 'var(--destructive)' }} />
                      <span className="truncate">{line}</span>
                    </div>
                    <button className="btn-icon" onClick={() => removeExcludePath(index)} title="Remove">
                      <X size={16} />
                    </button>
                  </div>
                ))}
                <button onClick={() => openBrowser('exclude')} className="btn-secondary flex items-center gap-2 w-full">
                  <Plus size={16} />
                  Add Exclude Path...
                </button>
              </div>
              <p className="text-xs mt-1" style={{ color: 'var(--text-muted)' }}>Directory names to skip during scanning.</p>
            </div>

            {/* Max Depth */}
            <div>
              <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-primary)' }}>Max Depth</label>
              <input
                type="number"
                value={maxDepth}
                onChange={(e) => setMaxDepth(parseInt(e.target.value, 10))}
                min="1"
                max="20"
                className="input-field w-32"
              />
              <p className="text-xs mt-1" style={{ color: 'var(--text-muted)' }}>How deep to recurse into scan paths.</p>
            </div>

            {/* Default Output */}
            <div>
              <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-primary)' }}>Default Output</label>
              <select
                value={defaultOutput}
                onChange={(e) => setDefaultOutput(e.target.value)}
                className="input-field input-field-select w-40"
              >
                <option value="table">Table</option>
                <option value="json">JSON</option>
                <option value="markdown">Markdown</option>
              </select>
            </div>

            {/* Actions */}
            <div className="flex items-center gap-4">
              <button onClick={saveConfig} disabled={saving} className="btn-primary flex items-center gap-2">
                {saving ? (
                  <>
                    <RefreshCw size={16} className="animate-spin" />
                    Saving...
                  </>
                ) : (
                  <>
                    <Save size={16} />
                    Save Settings
                  </>
                )}
              </button>
              {success && (
                <span className="flex items-center gap-1 text-sm" style={{ color: 'var(--success)' }}>
                  <CheckCircle2 size={14} /> Saved!
                </span>
              )}
            </div>

            {error && (
              <div className="card p-4 border-[var(--destructive)]">
                <p className="text-sm" style={{ color: 'var(--destructive)' }}>{error}</p>
              </div>
            )}
          </div>
        )}
      </div>

      {/* File Browser Modal */}
      {showBrowser && (
        <FileSystemBrowser
          selectFolder={true}
          onSelect={handlePathSelect}
          onClose={() => setShowBrowser(false)}
        />
      )}
    </div>
  );
}
