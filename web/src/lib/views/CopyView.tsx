import { useState } from 'react';
import { Copy as CopyIcon, RefreshCw, ArrowRight, CheckCircle2, XCircle } from 'lucide-react';
import { api, CopyResponse } from '../api-client';
import PathInput from '../components/PathInput';

type EntityType = 'skills' | 'hooks' | 'agents' | 'commands';

const entityTypes: EntityType[] = ['skills', 'hooks', 'agents', 'commands'];

export default function CopyView() {
  const [entityType, setEntityType] = useState<EntityType>('skills');
  const [fromPath, setFromPath] = useState('global');
  const [toPath, setToPath] = useState('');
  const [entityNames, setEntityNames] = useState('all');
  const [force, setForce] = useState(false);
  const [dryRun, setDryRun] = useState(true);
  const [result, setResult] = useState<CopyResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function runCopy() {
    if (!fromPath || !toPath) {
      setError('Please specify both source and destination');
      return;
    }
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const names = entityNames === 'all'
        ? ['all']
        : entityNames.split(',').map(n => n.trim());
      const data = await api.copy({
        type: entityType,
        from: fromPath,
        to: toPath,
        names,
        force,
        dry_run: dryRun,
      });
      setResult(data);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  function getStatusIcon(status: string) {
    if (status === 'copied' || status === 'would copy') return CheckCircle2;
    return XCircle;
  }

  return (
    <div className="space-y-6">
      {/* Copy Controls */}
      <div className="card p-6">
        <h3 className="text-lg font-semibold mb-4" style={{ color: 'var(--text-primary)' }}>Copy Entity</h3>
        <div className="space-y-4">
          <div className="flex flex-col md:flex-row gap-4 items-end">
            <div className="w-32">
              <label className="block text-sm mb-1" style={{ color: 'var(--text-secondary)' }}>Type</label>
              <select
                value={entityType}
                onChange={(e) => setEntityType(e.target.value as EntityType)}
                className="input-field input-field-select"
              >
                {entityTypes.map((type) => (
                  <option key={type} value={type}>{type}</option>
                ))}
              </select>
            </div>
            <div className="flex-1 w-full">
              <PathInput value={fromPath} placeholder="global or /path" onChange={setFromPath} />
            </div>
            <div className="flex items-center justify-center py-2">
              <ArrowRight size={20} style={{ color: 'var(--text-muted)' }} />
            </div>
            <div className="flex-1 w-full">
              <PathInput value={toPath} placeholder="/path/to/project" allowGlobal={false} onChange={setToPath} />
            </div>
          </div>
          <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
            <div className="flex-1 w-full">
              <label className="block text-sm mb-1" style={{ color: 'var(--text-secondary)' }}>
                Entity names (comma-separated, or "all")
              </label>
              <input
                type="text"
                value={entityNames}
                onChange={(e) => setEntityNames(e.target.value)}
                placeholder="all"
                className="input-field"
              />
            </div>
            <div className="flex items-center gap-4">
              <label className="flex items-center gap-2 text-sm" style={{ color: 'var(--text-secondary)' }}>
                <input
                  type="checkbox"
                  checked={force}
                  onChange={(e) => setForce(e.target.checked)}
                />
                Force overwrite
              </label>
              <label className="flex items-center gap-2 text-sm" style={{ color: 'var(--text-secondary)' }}>
                <input
                  type="checkbox"
                  checked={dryRun}
                  onChange={(e) => setDryRun(e.target.checked)}
                />
                Dry run
              </label>
            </div>
          </div>
          <button
            onClick={runCopy}
            disabled={loading}
            className="btn-primary flex items-center gap-2 whitespace-nowrap w-full md:w-auto"
          >
            {loading ? (
              <RefreshCw size={16} className="animate-spin" />
            ) : (
              <CopyIcon size={16} />
            )}
            {dryRun ? 'Preview Copy' : 'Copy Entity'}
          </button>
        </div>
      </div>

      {error && (
        <div className="card p-4 border-[var(--destructive)]">
          <p className="text-sm" style={{ color: 'var(--destructive)' }}>{error}</p>
        </div>
      )}

      {result && (
        <div className="card p-4">
          <h4 className="text-sm font-semibold mb-3" style={{ color: 'var(--text-primary)' }}>Results</h4>
          <div className="space-y-1">
            {result.results?.map((r) => {
              const Icon = getStatusIcon(r.status);
              return (
                <div key={r.name} className="flex items-center gap-2 py-1">
                  <Icon
                    size={14}
                    style={{ color: r.status === 'copied' || r.status === 'would copy' ? 'var(--success)' : 'var(--destructive)' }}
                  />
                  <span className="font-mono text-sm" style={{ color: 'var(--text-primary)' }}>{r.name}</span>
                  {r.detail && (
                    <span className="text-xs" style={{ color: 'var(--text-muted)' }}>— {r.detail}</span>
                  )}
                </div>
              );
            })}
          </div>
          {result.message && (
            <p className="mt-3 text-sm" style={{ color: 'var(--text-secondary)' }}>{result.message}</p>
          )}
        </div>
      )}
    </div>
  );
}
