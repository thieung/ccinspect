import { useState } from 'react';
import { GitCompare, RefreshCw, Plus, Minus, ArrowRight } from 'lucide-react';
import { api, DiffResponse } from '../api-client';
import PathInput from '../components/PathInput';

type EntityType = 'skills' | 'hooks' | 'agents' | 'commands';

export default function DiffView() {
  const [pathA, setPathA] = useState('global');
  const [pathB, setPathB] = useState('');
  const [entityType, setEntityType] = useState<EntityType>('skills');
  const [diffData, setDiffData] = useState<DiffResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function runDiff() {
    if (!pathA || !pathB) {
      setError('Please select two projects to compare');
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const data = await api.diff(pathA, pathB, entityType);
      setDiffData(data);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      {/* Diff Controls */}
      <div className="card p-6">
        <h3 className="text-lg font-semibold mb-4" style={{ color: 'var(--text-primary)' }}>Compare Projects</h3>
        <div className="flex flex-col md:flex-row gap-4 items-end">
          <div className="flex-1 w-full">
            <PathInput
              value={pathA}
              placeholder="global or /path/to/project"
              onChange={setPathA}
            />
          </div>
          <div className="flex items-center justify-center py-2">
            <ArrowRight size={20} style={{ color: 'var(--text-muted)' }} />
          </div>
          <div className="flex-1 w-full">
            <PathInput
              value={pathB}
              placeholder="/path/to/project"
              allowGlobal={false}
              onChange={setPathB}
            />
          </div>
          <div className="w-32">
            <label className="block text-sm mb-1" style={{ color: 'var(--text-secondary)' }}>Type</label>
            <select
              value={entityType}
              onChange={(e) => setEntityType(e.target.value as EntityType)}
              className="input-field input-field-select"
            >
              <option value="skills">Skills</option>
              <option value="hooks">Hooks</option>
              <option value="agents">Agents</option>
              <option value="commands">Commands</option>
            </select>
          </div>
          <button onClick={runDiff} disabled={loading} className="btn-primary flex items-center gap-2 whitespace-nowrap">
            {loading ? (
              <RefreshCw size={16} className="animate-spin" />
            ) : (
              <GitCompare size={16} />
            )}
            Compare
          </button>
        </div>
      </div>

      {error && (
        <div className="card p-4 border-[var(--destructive)]">
          <p className="text-sm" style={{ color: 'var(--destructive)' }}>{error}</p>
        </div>
      )}

      {diffData && (
        <>
          {/* Results Header */}
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold" style={{ color: 'var(--text-primary)' }}>
              Diff: {diffData.left_name} vs {diffData.right_name}
            </h3>
            <div className="flex gap-4 text-sm">
              <span className="flex items-center gap-1" style={{ color: 'var(--success)' }}>
                <Plus size={14} /> {diffData.added.length} added
              </span>
              <span className="flex items-center gap-1" style={{ color: 'var(--destructive)' }}>
                <Minus size={14} /> {diffData.removed.length} removed
              </span>
            </div>
          </div>

          {/* Added (in B but not in A) */}
          {diffData.added.length > 0 && (
            <div className="card p-4">
              <h4 className="text-sm font-semibold mb-3" style={{ color: 'var(--success)' }}>
                Added in {diffData.right_name}
              </h4>
              <div className="space-y-1">
                {diffData.added.map((item) => (
                  <div key={item.name} className="flex items-center gap-2 py-1">
                    <Plus size={14} style={{ color: 'var(--success)' }} />
                    <span className="font-mono text-sm" style={{ color: 'var(--text-primary)' }}>{item.name}</span>
                    {item.description && (
                      <span className="text-xs" style={{ color: 'var(--text-muted)' }}>— {item.description}</span>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Removed (in A but not in B) */}
          {diffData.removed.length > 0 && (
            <div className="card p-4">
              <h4 className="text-sm font-semibold mb-3" style={{ color: 'var(--destructive)' }}>
                Removed from {diffData.left_name}
              </h4>
              <div className="space-y-1">
                {diffData.removed.map((item) => (
                  <div key={item.name} className="flex items-center gap-2 py-1">
                    <Minus size={14} style={{ color: 'var(--destructive)' }} />
                    <span className="font-mono text-sm" style={{ color: 'var(--text-primary)' }}>{item.name}</span>
                    {item.description && (
                      <span className="text-xs" style={{ color: 'var(--text-muted)' }}>— {item.description}</span>
                    )}
                  </div>
                ))}
              </div>
            </div>
          )}

          {diffData.added.length === 0 && diffData.removed.length === 0 && (
            <div className="card p-8 text-center" style={{ color: 'var(--text-muted)' }}>
              <p>No differences found — both projects have the same {entityType}.</p>
            </div>
          )}
        </>
      )}
    </div>
  );
}
