import { useState, useEffect } from 'react';
import { api, ListResponse } from '../api-client';
import {
  Package, Webhook, Bot, FileText, Server, Users, Copy, ChevronRight, RefreshCw
} from 'lucide-react';

type EntityType = 'skills' | 'hooks' | 'agents' | 'commands' | 'rules' | 'mcp' | 'teams';

const entityFilters: EntityType[] = ['skills', 'hooks', 'agents', 'commands', 'rules', 'mcp', 'teams'];

const entityIcons: Record<string, React.ComponentType<{ size: number }>> = {
  skills: Package,
  hooks: Webhook,
  agents: Bot,
  commands: FileText,
  rules: FileText,
  mcp: Server,
  teams: Users,
};

export default function EntitiesView() {
  const [entityType, setEntityType] = useState<EntityType>('skills');
  const [prefix, setPrefix] = useState('');
  const [globalOnly, setGlobalOnly] = useState(false);
  const [listData, setListData] = useState<ListResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function loadEntities() {
    setLoading(true);
    setError(null);
    try {
      const data = await api.list(entityType, { global: globalOnly, prefix });
      setListData(data);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadEntities();
  }, [entityType, globalOnly]);

  function getIcon(type: string) {
    return entityIcons[type] || Package;
  }

  function shortenPath(path: string): string {
    if (path === 'global') return 'global';
    const parts = path.split('/');
    return '~/' + parts.slice(-2).join('/');
  }

  const handlePrefixKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      loadEntities();
    }
  };

  return (
    <div className="space-y-6">
      {/* Filter Bar */}
      <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
        <div className="flex items-center gap-2 overflow-x-auto pb-2">
          {entityFilters.map((filter) => (
            <button
              key={filter}
              onClick={() => { setEntityType(filter); loadEntities(); }}
              className="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-all duration-200"
              style={{
                backgroundColor: entityType === filter ? 'var(--accent)' : 'var(--bg-muted)',
                color: entityType === filter ? 'white' : 'var(--text-secondary)'
              }}
            >
              {filter.charAt(0).toUpperCase() + filter.slice(1)}
            </button>
          ))}
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <input
            type="text"
            value={prefix}
            onChange={(e) => setPrefix(e.target.value)}
            onKeyDown={handlePrefixKeyDown}
            placeholder="Prefix filter..."
            className="input-field w-40 text-sm"
          />
          <label className="flex items-center gap-2 text-sm whitespace-nowrap" style={{ color: 'var(--text-secondary)' }}>
            <input
              type="checkbox"
              checked={globalOnly}
              onChange={(e) => setGlobalOnly(e.target.checked)}
            />
            Global only
          </label>
          <button onClick={loadEntities} disabled={loading} className="btn-secondary text-sm">
            <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
          </button>
        </div>
      </div>

      {error && (
        <div className="card p-4 border-[var(--destructive)]">
          <p className="text-sm" style={{ color: 'var(--destructive)' }}>Error: {error}</p>
        </div>
      )}

      {loading && !listData ? (
        <div className="text-center py-12" style={{ color: 'var(--text-muted)' }}>
          <RefreshCw size={32} className="animate-spin mx-auto mb-4" />
          <p>Loading entities...</p>
        </div>
      ) : (
        /* Entities List */
        <div className="space-y-2">
          {listData?.entities.map((entity, i) => {
            const Icon = getIcon(entity.type ?? entityType);
            return (
              <div key={`${entity.name}-${entity.source}`} className="entity-list-item animate-in" style={{ animationDelay: `${i * 0.03}s` }}>
                <div className="entity-list-item-icon">
                  <Icon size={18} />
                </div>
                <div className="entity-list-item-content">
                  <p className="entity-list-item-title">{entity.name}</p>
                  {entity.description && (
                    <p className="entity-list-item-meta">{entity.description}</p>
                  )}
                </div>
                <span className={`badge ${entity.source === 'global' ? 'badge-brand' : 'badge-success'}`}>
                  {shortenPath(entity.source)}
                </span>
                <div className="flex items-center gap-1">
                  <button className="btn-icon w-8 h-8" aria-label="Copy" title="Copy entity">
                    <Copy size={14} />
                  </button>
                  <button className="btn-icon w-8 h-8" aria-label="View details">
                    <ChevronRight size={14} />
                  </button>
                </div>
              </div>
            );
          })}
          {listData?.entities.length === 0 && (
            <div className="text-center py-12" style={{ color: 'var(--text-muted)' }}>
              <p>No {entityType} found</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
