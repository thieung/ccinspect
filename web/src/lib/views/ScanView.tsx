import { useState, useEffect, useMemo } from 'react';
import {
  Search, Scan, LayoutGrid, List as ListIcon, Terminal, Package, Webhook, Bot, Server, Activity,
  CheckCircle2, RefreshCw, ChevronRight, FolderSearch
} from 'lucide-react';
import { api, ScanResponse } from '../api-client';
import ProjectDetailModal from '../components/ProjectDetailModal';

interface ScanViewProps {
  onNavigate?: (tab: string, projectPath?: string) => void;
}

function getProjectName(path: string): string {
  const parts = path.split('/');
  return parts[parts.length - 1] || path;
}

export default function ScanView({ onNavigate }: ScanViewProps) {
  const [scanData, setScanData] = useState<ScanResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedProject, setSelectedProject] = useState<ScanResponse['projects'][number] | null>(null);

  async function loadScan() {
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

  const filteredProjects = useMemo(() => {
    return scanData?.projects.filter(p =>
      getProjectName(p.path).toLowerCase().includes(searchQuery.toLowerCase())
    ) ?? [];
  }, [scanData?.projects, searchQuery]);

  const totalSkills = useMemo(() =>
    scanData?.projects.reduce((s, p) => s + p.skills, 0) ?? 0
  , [scanData?.projects]);

  const totalHooks = useMemo(() =>
    scanData?.projects.reduce((s, p) => s + p.hooks, 0) ?? 0
  , [scanData?.projects]);

  const totalAgents = useMemo(() =>
    scanData?.projects.reduce((s, p) => s + p.agents, 0) ?? 0
  , [scanData?.projects]);

  const totalMcp = useMemo(() =>
    scanData?.projects.reduce((s, p) => s + p.mcp_servers, 0) ?? 0
  , [scanData?.projects]);

  useEffect(() => {
    loadScan();
  }, []);

  return (
    <div className="space-y-8">
      {/* Stats Grid */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
        <div className="stat-card animate-in stagger-1">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ backgroundColor: 'var(--accent-muted)' }}>
              <FolderSearch size={16} style={{ color: 'var(--accent)' }} />
            </div>
            <span className="text-xs" style={{ color: 'var(--text-muted)' }}>Projects</span>
          </div>
          <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{scanData?.projects.length ?? 0}</p>
        </div>
        <div className="stat-card animate-in stagger-2">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ backgroundColor: 'var(--accent-muted)' }}>
              <Package size={16} style={{ color: 'var(--accent)' }} />
            </div>
            <span className="text-xs" style={{ color: 'var(--text-muted)' }}>Skills</span>
          </div>
          <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{totalSkills}</p>
        </div>
        <div className="stat-card animate-in stagger-3">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ backgroundColor: 'var(--accent-muted)' }}>
              <Webhook size={16} style={{ color: 'var(--accent)' }} />
            </div>
            <span className="text-xs" style={{ color: 'var(--text-muted)' }}>Hooks</span>
          </div>
          <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{totalHooks}</p>
        </div>
        <div className="stat-card animate-in stagger-4">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ backgroundColor: 'var(--accent-muted)' }}>
              <Bot size={16} style={{ color: 'var(--accent)' }} />
            </div>
            <span className="text-xs" style={{ color: 'var(--text-muted)' }}>Agents</span>
          </div>
          <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{totalAgents}</p>
        </div>
        <div className="stat-card animate-in stagger-5">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ backgroundColor: 'var(--accent-muted)' }}>
              <Server size={16} style={{ color: 'var(--accent)' }} />
            </div>
            <span className="text-xs" style={{ color: 'var(--text-muted)' }}>MCP</span>
          </div>
          <p className="text-2xl font-bold" style={{ color: 'var(--text-primary)' }}>{totalMcp}</p>
        </div>
        <div className="stat-card animate-in stagger-6">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ backgroundColor: 'var(--accent-muted)' }}>
              <Activity size={16} style={{ color: 'var(--accent)' }} />
            </div>
            <span className="text-xs" style={{ color: 'var(--text-muted)' }}>Status</span>
          </div>
          <p className="text-lg font-bold" style={{ color: 'var(--success)' }}>{loading ? 'Scanning' : 'Active'}</p>
        </div>
      </div>

      {error && (
        <div className="card p-4 border-[var(--destructive)]">
          <p className="text-sm" style={{ color: 'var(--destructive)' }}>Error: {error}</p>
          <button onClick={loadScan} className="btn-secondary mt-2 text-sm">Retry</button>
        </div>
      )}

      {/* Action Bar */}
      <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
        <div className="relative w-full sm:w-80">
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search projects..."
            className="input-field pl-11"
          />
          <Search size={18} className="absolute left-3 top-1/2 -translate-y-1/2 z-10" style={{ color: 'var(--text-muted)' }} />
        </div>
        <div className="flex items-center gap-2">
          <button onClick={loadScan} disabled={loading} className="btn-primary flex items-center gap-2 whitespace-nowrap">
            {loading ? (
              <>
                <RefreshCw size={16} className="animate-spin" />
                Scanning...
              </>
            ) : (
              <>
                <Scan size={16} />
                Scan Projects
              </>
            )}
          </button>
          <div className="flex items-center gap-1 p-1 rounded-lg" style={{ backgroundColor: 'var(--bg-muted)' }}>
            <button
              onClick={() => setViewMode('grid')}
              className="p-2 rounded-md transition-all duration-200"
              style={{
                backgroundColor: viewMode === 'grid' ? 'var(--bg-card)' : 'transparent',
                color: viewMode === 'grid' ? 'var(--accent)' : 'var(--text-muted)'
              }}
              aria-label="Grid view"
            >
              <LayoutGrid size={18} />
            </button>
            <button
              onClick={() => setViewMode('list')}
              className="p-2 rounded-md transition-all duration-200"
              style={{
                backgroundColor: viewMode === 'list' ? 'var(--bg-card)' : 'transparent',
                color: viewMode === 'list' ? 'var(--accent)' : 'var(--text-muted)'
              }}
              aria-label="List view"
            >
              <ListIcon size={18} />
            </button>
          </div>
        </div>
      </div>

      {/* Projects Grid/List */}
      {loading && !scanData ? (
        <div className="text-center py-12" style={{ color: 'var(--text-muted)' }}>
          <RefreshCw size={32} className="animate-spin mx-auto mb-4" />
          <p>Scanning projects...</p>
        </div>
      ) : viewMode === 'grid' ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredProjects.map((project, i) => (
            <div
              key={project.path}
              className="card p-5 cursor-pointer animate-in hover:border-[var(--accent)] transition-all"
              style={{ animationDelay: `${i * 0.05}s` }}
              onClick={() => setSelectedProject(project)}
            >
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'var(--accent)' }}>
                    <Terminal size={20} className="text-white" />
                  </div>
                  <div>
                    <h3 className="font-semibold" style={{ color: 'var(--text-primary)' }}>{getProjectName(project.path)}</h3>
                    <p className="text-xs font-mono truncate max-w-48" style={{ color: 'var(--text-muted)' }}>{project.path}</p>
                  </div>
                </div>
                <button className="btn-icon w-8 h-8" aria-label="Project menu">
                  <ChevronRight size={16} />
                </button>
              </div>
              <div className="grid grid-cols-4 gap-2 mb-3">
                <div className="text-center p-2 rounded-lg" style={{ backgroundColor: 'var(--bg-muted)' }}>
                  <p className="text-lg font-bold" style={{ color: 'var(--accent)' }}>{project.skills}</p>
                  <p className="text-xs" style={{ color: 'var(--text-muted)' }}>Skills</p>
                </div>
                <div className="text-center p-2 rounded-lg" style={{ backgroundColor: 'var(--bg-muted)' }}>
                  <p className="text-lg font-bold" style={{ color: 'var(--accent)' }}>{project.hooks}</p>
                  <p className="text-xs" style={{ color: 'var(--text-muted)' }}>Hooks</p>
                </div>
                <div className="text-center p-2 rounded-lg" style={{ backgroundColor: 'var(--bg-muted)' }}>
                  <p className="text-lg font-bold" style={{ color: 'var(--accent)' }}>{project.agents}</p>
                  <p className="text-xs" style={{ color: 'var(--text-muted)' }}>Agents</p>
                </div>
                <div className="text-center p-2 rounded-lg" style={{ backgroundColor: 'var(--bg-muted)' }}>
                  <p className="text-lg font-bold" style={{ color: 'var(--accent)' }}>{project.mcp_servers}</p>
                  <p className="text-xs" style={{ color: 'var(--text-muted)' }}>MCP</p>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className={`badge ${project.has_claude_md ? 'badge-success' : 'badge-warning'}`}>
                  <CheckCircle2 size={12} className="mr-1" />
                  {project.has_claude_md ? 'Has CLAUDE.md' : 'No CLAUDE.md'}
                </span>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="space-y-2">
          {filteredProjects.map((project, i) => (
            <div
              key={project.path}
              className="entity-list-item animate-in cursor-pointer hover:border-[var(--accent)] transition-all"
              style={{ animationDelay: `${i * 0.03}s` }}
              onClick={() => setSelectedProject(project)}
            >
              <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'var(--accent)' }}>
                <Terminal size={18} className="text-white" />
              </div>
              <div className="entity-list-item-content">
                <p className="entity-list-item-title">{getProjectName(project.path)}</p>
                <p className="entity-list-item-meta font-mono text-xs">{project.path}</p>
              </div>
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-3 text-sm">
                  <span className="flex items-center gap-1" style={{ color: 'var(--text-muted)' }}>
                    <Package size={14} /> {project.skills}
                  </span>
                  <span className="flex items-center gap-1" style={{ color: 'var(--text-muted)' }}>
                    <Webhook size={14} /> {project.hooks}
                  </span>
                  <span className="flex items-center gap-1" style={{ color: 'var(--text-muted)' }}>
                    <Bot size={14} /> {project.agents}
                  </span>
                </div>
                <span className={`badge ${project.has_claude_md ? 'badge-success' : 'badge-warning'}`}>
                  {project.has_claude_md ? 'CLAUDE.md' : 'No CLAUDE.md'}
                </span>
              </div>
              <ChevronRight size={18} style={{ color: 'var(--text-muted)' }} />
            </div>
          ))}
        </div>
      )}

      {selectedProject && (
        <ProjectDetailModal
          project={selectedProject}
          onClose={() => setSelectedProject(null)}
          onNavigate={(tab) => {
            onNavigate?.(tab);
            setSelectedProject(null);
          }}
        />
      )}
    </div>
  );
}
