import { List, GitCompare, Copy, ExternalLink, Terminal, CheckCircle2 } from 'lucide-react';
import { ProjectSummary } from '../api-client';

interface ProjectDetailModalProps {
  project: ProjectSummary;
  onClose?: () => void;
  onNavigate?: (tab: string) => void;
}

function getProjectName(path: string): string {
  const parts = path.split('/');
  return parts[parts.length - 1] || path;
}

export default function ProjectDetailModal({ project, onClose, onNavigate }: ProjectDetailModalProps) {
  function closeModal() {
    onClose?.();
  }

  function handleNavigate(tab: string) {
    onNavigate?.(tab);
    closeModal();
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      style={{ backgroundColor: 'rgba(0, 0, 0, 0.5)' }}
      onClick={() => closeModal()}
    >
      <div
        className="card w-full max-w-lg"
        style={{ backgroundColor: 'var(--bg-card)' }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b" style={{ border: '1px solid var(--border-color)' }}>
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'var(--accent)' }}>
              <Terminal size={20} className="text-white" />
            </div>
            <div>
              <h3 className="font-semibold" style={{ color: 'var(--text-primary)' }}>
                {getProjectName(project.path)}
              </h3>
              <p className="text-xs font-mono truncate max-w-48" style={{ color: 'var(--text-muted)' }}>
                {project.path}
              </p>
            </div>
          </div>
          <button
            type="button"
            className="btn-icon w-8 h-8"
            onClick={() => closeModal()}
            title="Close (Esc)"
          >
            <span style={{ fontSize: '20px' }}>×</span>
          </button>
        </div>

        {/* Stats */}
        <div className="p-4">
          <div className="grid grid-cols-4 gap-2 mb-4">
            <div className="text-center p-3 rounded-lg" style={{ backgroundColor: 'var(--bg-muted)' }}>
              <p className="text-2xl font-bold" style={{ color: 'var(--accent)' }}>{project.skills}</p>
              <p className="text-xs" style={{ color: 'var(--text-muted)' }}>Skills</p>
            </div>
            <div className="text-center p-3 rounded-lg" style={{ backgroundColor: 'var(--bg-muted)' }}>
              <p className="text-2xl font-bold" style={{ color: 'var(--accent)' }}>{project.hooks}</p>
              <p className="text-xs" style={{ color: 'var(--text-muted)' }}>Hooks</p>
            </div>
            <div className="text-center p-3 rounded-lg" style={{ backgroundColor: 'var(--bg-muted)' }}>
              <p className="text-2xl font-bold" style={{ color: 'var(--accent)' }}>{project.agents}</p>
              <p className="text-xs" style={{ color: 'var(--text-muted)' }}>Agents</p>
            </div>
            <div className="text-center p-3 rounded-lg" style={{ backgroundColor: 'var(--bg-muted)' }}>
              <p className="text-2xl font-bold" style={{ color: 'var(--accent)' }}>{project.mcp_servers}</p>
              <p className="text-xs" style={{ color: 'var(--text-muted)' }}>MCP</p>
            </div>
          </div>

          {/* CLAUDE.md Status */}
          <div
            className="flex items-center gap-2 p-3 rounded-lg mb-4"
            style={{
              backgroundColor: project.has_claude_md ? 'var(--success-muted)' : 'var(--warning-muted)',
              color: project.has_claude_md ? 'var(--success)' : 'var(--warning)'
            }}
          >
            <CheckCircle2 size={16} style={{ color: project.has_claude_md ? 'var(--success)' : 'var(--warning)' }} />
            <span className="text-sm font-medium">
              {project.has_claude_md ? 'Has CLAUDE.md' : 'No CLAUDE.md'}
            </span>
          </div>

          {/* Actions */}
          <div className="space-y-2">
            <button
              onClick={() => handleNavigate('entities')}
              className="w-full flex items-center gap-3 p-3 rounded-lg hover:border-[var(--accent)] transition-all"
              style={{ backgroundColor: 'var(--bg-muted)', border: '1px solid var(--border-subtle)' }}
            >
              <List size={18} style={{ color: 'var(--accent)' }} />
              <span className="font-medium" style={{ color: 'var(--text-primary)' }}>View Entities</span>
              <ExternalLink size={14} style={{ color: 'var(--text-muted)', marginLeft: 'auto' }} />
            </button>
            <button
              onClick={() => handleNavigate('diff')}
              className="w-full flex items-center gap-3 p-3 rounded-lg hover:border-[var(--accent)] transition-all"
              style={{ backgroundColor: 'var(--bg-muted)', border: '1px solid var(--border-subtle)' }}
            >
              <GitCompare size={18} style={{ color: 'var(--accent)' }} />
              <span className="font-medium" style={{ color: 'var(--text-primary)' }}>Compare with Global</span>
              <ExternalLink size={14} style={{ color: 'var(--text-muted)', marginLeft: 'auto' }} />
            </button>
            <button
              onClick={() => handleNavigate('copy')}
              className="w-full flex items-center gap-3 p-3 rounded-lg hover:border-[var(--accent)] transition-all"
              style={{ backgroundColor: 'var(--bg-muted)', border: '1px solid var(--border-subtle)' }}
            >
              <Copy size={18} style={{ color: 'var(--accent)' }} />
              <span className="font-medium" style={{ color: 'var(--text-primary)' }}>Copy from Global</span>
              <ExternalLink size={14} style={{ color: 'var(--text-muted)', marginLeft: 'auto' }} />
            </button>
          </div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-2 p-4 border-t" style={{ borderColor: 'var(--border-color)' }}>
          <button type="button" className="btn-secondary" onClick={() => closeModal()}>Close</button>
        </div>
      </div>
    </div>
  );
}
