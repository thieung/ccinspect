import { useState, useEffect } from 'react';
import { api, ScanResponse } from '../api-client';
import { Globe, Package, Webhook, Bot, Copy, GitCompare, Trash2, RefreshCw, FileText, Users } from 'lucide-react';

interface GlobalViewProps {
  onNavigate?: (tab: string) => void;
}

type SectionType = 'skills' | 'hooks' | 'agents' | 'commands' | 'rules' | 'teams' | null;

export default function GlobalView({ onNavigate }: GlobalViewProps) {
  const [scanData, setScanData] = useState<ScanResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [expandedSection, setExpandedSection] = useState<SectionType>(null);

  async function loadGlobal() {
    setLoading(true);
    setError(null);
    try {
      const data = await api.scan();
      setScanData(data);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadGlobal();
  }, []);

  function toggleSection(section: SectionType) {
    setExpandedSection(expandedSection === section ? null : section);
  }

  if (loading) {
    return (
      <div className="text-center py-12" style={{ color: 'var(--text-muted)' }}>
        <RefreshCw size={32} className="animate-spin mx-auto mb-4" />
        <p>Loading global config...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="card p-4 border-[var(--destructive)]">
        <p className="text-sm" style={{ color: 'var(--destructive)' }}>Error: {error}</p>
        <button onClick={loadGlobal} className="btn-secondary mt-2 text-sm">Retry</button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="card p-6">
        <div className="flex items-center gap-3 mb-6">
          <div className="w-12 h-12 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'var(--accent)' }}>
            <Globe size={24} className="text-white" />
          </div>
          <div>
            <h2 className="text-xl font-semibold" style={{ color: 'var(--text-primary)' }}>Global Configuration</h2>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>{scanData?.global?.path ?? '~/.claude/'}</p>
          </div>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div
            className="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
            style={{ backgroundColor: 'var(--bg-muted)', borderColor: 'var(--border-color)' }}
            onClick={() => toggleSection('skills')}
          >
            <Package size={20} style={{ color: 'var(--accent)' }} className="mb-2" />
            <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{scanData?.global?.skills ?? 0}</p>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Skills</p>
          </div>
          <div
            className="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
            style={{ backgroundColor: 'var(--bg-muted)', borderColor: 'var(--border-color)' }}
            onClick={() => toggleSection('hooks')}
          >
            <Webhook size={20} style={{ color: 'var(--accent)' }} className="mb-2" />
            <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{scanData?.global?.hooks ?? 0}</p>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Hooks</p>
          </div>
          <div
            className="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
            style={{ backgroundColor: 'var(--bg-muted)', borderColor: 'var(--border-color)' }}
            onClick={() => toggleSection('agents')}
          >
            <Bot size={20} style={{ color: 'var(--accent)' }} className="mb-2" />
            <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{scanData?.global?.agents ?? 0}</p>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Agents</p>
          </div>
          <div
            className="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
            style={{ backgroundColor: 'var(--bg-muted)', borderColor: 'var(--border-color)' }}
            onClick={() => toggleSection('commands')}
          >
            <Copy size={20} style={{ color: 'var(--accent)' }} className="mb-2" />
            <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{scanData?.global?.commands ?? 0}</p>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Commands</p>
          </div>
          <div
            className="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
            style={{ backgroundColor: 'var(--bg-muted)', borderColor: 'var(--border-color)' }}
            onClick={() => toggleSection('rules')}
          >
            <FileText size={20} style={{ color: 'var(--accent)' }} className="mb-2" />
            <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{scanData?.global?.rules ?? 0}</p>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Rules</p>
          </div>
          <div
            className="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
            style={{ backgroundColor: 'var(--bg-muted)', borderColor: 'var(--border-color)' }}
            onClick={() => toggleSection('teams')}
          >
            <Users size={20} style={{ color: 'var(--accent)' }} className="mb-2" />
            <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{scanData?.global?.teams ?? 0}</p>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Teams</p>
          </div>
        </div>

        {/* Expanded Sections */}
        {expandedSection && (
          <div className="mt-6 p-4 rounded-xl border" style={{ backgroundColor: 'var(--bg-card)', borderColor: 'var(--border-color)' }}>
            <div className="flex items-center justify-between mb-3">
              <h3 className="font-semibold" style={{ color: 'var(--text-primary)' }}>
                {expandedSection.charAt(0).toUpperCase() + expandedSection.slice(1)}
              </h3>
              <button className="btn-icon w-6 h-6" onClick={() => setExpandedSection(null)}>
                <span style={{ fontSize: '16px' }}>×</span>
              </button>
            </div>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>
              Use the <strong>Entities</strong> tab to view and manage {expandedSection} in detail.
            </p>
            <button onClick={() => onNavigate?.('entities')} className="btn-primary mt-3 text-sm">
              View All {expandedSection.charAt(0).toUpperCase() + expandedSection.slice(1)}
            </button>
          </div>
        )}
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <button
          onClick={() => onNavigate?.('copy')}
          className="card p-5 flex items-center gap-4 hover:border-[var(--accent)] transition-all duration-200 cursor-pointer"
        >
          <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'var(--accent-muted)' }}>
            <Copy size={20} style={{ color: 'var(--accent)' }} />
          </div>
          <div className="text-left">
            <p className="font-medium" style={{ color: 'var(--text-primary)' }}>Copy Entity</p>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Copy to project</p>
          </div>
        </button>
        <button
          onClick={() => onNavigate?.('diff')}
          className="card p-5 flex items-center gap-4 hover:border-[var(--accent)] transition-all duration-200 cursor-pointer"
        >
          <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'var(--accent-muted)' }}>
            <GitCompare size={20} style={{ color: 'var(--accent)' }} />
          </div>
          <div className="text-left">
            <p className="font-medium" style={{ color: 'var(--text-primary)' }}>Diff Skills</p>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Compare projects</p>
          </div>
        </button>
        <button
          onClick={() => onNavigate?.('clean')}
          className="card p-5 flex items-center gap-4 hover:border-[var(--accent)] transition-all duration-200 cursor-pointer"
        >
          <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'var(--accent-muted)' }}>
            <Trash2 size={20} style={{ color: 'var(--destructive)' }} />
          </div>
          <div className="text-left">
            <p className="font-medium" style={{ color: 'var(--text-primary)' }}>Clean Teams</p>
            <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Remove stale teams</p>
          </div>
        </button>
      </div>
    </div>
  );
}
