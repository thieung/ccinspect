import { useState } from 'react';
import { Trash2, RefreshCw, AlertTriangle, CheckCircle2 } from 'lucide-react';
import { api, CleanResponse, CleanTeamsResponse } from '../api-client';
import PathInput from '../components/PathInput';

type CleanMode = 'project' | 'teams';

export default function CleanView() {
  const [mode, setMode] = useState<CleanMode>('project');

  // Project clean
  const [projectPath, setProjectPath] = useState('');
  const [projectDryRun, setProjectDryRun] = useState(true);
  const [projectResult, setProjectResult] = useState<CleanResponse | null>(null);

  // Teams clean
  const [teamsDryRun, setTeamsDryRun] = useState(true);
  const [teamsAll, setTeamsAll] = useState(false);
  const [teamsResult, setTeamsResult] = useState<CleanTeamsResponse | null>(null);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function cleanProject() {
    if (!projectPath) {
      setError('Please specify a project path');
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const result = await api.clean({ path: projectPath, dry_run: projectDryRun });
      setProjectResult(result);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  async function cleanTeams() {
    setLoading(true);
    setError(null);
    try {
      const result = await api.cleanTeams({ all: teamsAll, dry_run: teamsDryRun });
      setTeamsResult(result);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      {/* Mode Toggle */}
      <div className="flex gap-2">
        <button
          onClick={() => setMode('project')}
          className="px-4 py-2 rounded-lg text-sm font-medium transition-all"
          style={{
            backgroundColor: mode === 'project' ? 'var(--accent)' : 'var(--bg-muted)',
            color: mode === 'project' ? 'white' : 'var(--text-secondary)'
          }}
        >
          Clean Project
        </button>
        <button
          onClick={() => setMode('teams')}
          className="px-4 py-2 rounded-lg text-sm font-medium transition-all"
          style={{
            backgroundColor: mode === 'teams' ? 'var(--accent)' : 'var(--bg-muted)',
            color: mode === 'teams' ? 'white' : 'var(--text-secondary)'
          }}
        >
          Clean Teams
        </button>
      </div>

      {mode === 'project' ? (
        <>
          {/* Project Clean */}
          <div className="card p-6">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'rgba(248,81,73,0.15)' }}>
                <Trash2 size={20} style={{ color: 'var(--destructive)' }} />
              </div>
              <div>
                <h3 className="text-lg font-semibold" style={{ color: 'var(--text-primary)' }}>Clean Project</h3>
                <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Remove .claude/ directory from a project</p>
              </div>
            </div>
            <div className="space-y-4">
              <PathInput
                value={projectPath}
                placeholder="/path/to/project"
                allowGlobal={false}
                label="Project Path"
                onChange={setProjectPath}
              />
              <label className="flex items-center gap-2 text-sm" style={{ color: 'var(--text-secondary)' }}>
                <input
                  type="checkbox"
                  checked={projectDryRun}
                  onChange={(e) => setProjectDryRun(e.target.checked)}
                />
                Dry run (preview only)
              </label>
              {projectDryRun && (
                <div className="flex items-center gap-2 text-sm" style={{ color: 'var(--warning)' }}>
                  <AlertTriangle size={14} />
                  No files will be deleted in dry run mode
                </div>
              )}
              <button
                onClick={cleanProject}
                disabled={loading}
                className="btn-primary flex items-center gap-2 whitespace-nowrap"
                style={{ backgroundColor: 'var(--destructive)' }}
              >
                {loading ? (
                  <RefreshCw size={16} className="animate-spin" />
                ) : (
                  <Trash2 size={16} />
                )}
                {projectDryRun ? 'Preview Clean' : 'Clean Project'}
              </button>
            </div>
          </div>

          {projectResult && (
            <div className="card p-4">
              <div className="flex items-center gap-2 mb-3">
                {projectResult.dry_run ? (
                  <>
                    <AlertTriangle size={16} style={{ color: 'var(--warning)' }} />
                    <h4 className="text-sm font-semibold" style={{ color: 'var(--warning)' }}>Dry Run Preview</h4>
                  </>
                ) : (
                  <>
                    <CheckCircle2 size={16} style={{ color: 'var(--success)' }} />
                    <h4 className="text-sm font-semibold" style={{ color: 'var(--success)' }}>Clean Complete</h4>
                  </>
                )}
              </div>
              <p className="text-sm mb-2" style={{ color: 'var(--text-primary)' }}>{projectResult.message}</p>
              <p className="text-sm" style={{ color: 'var(--text-muted)' }}>
                {projectResult.files_count} files would be/were removed
              </p>
              {projectResult.files && projectResult.files.length > 0 && (
                <div className="mt-3 max-h-48 overflow-y-auto">
                  {projectResult.files.map((file) => (
                    <p key={file} className="text-xs font-mono py-0.5" style={{ color: 'var(--text-secondary)' }}>{file}</p>
                  ))}
                </div>
              )}
            </div>
          )}
        </>
      ) : (
        <>
          {/* Teams Clean */}
          <div className="card p-6">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-10 h-10 rounded-xl flex items-center justify-center" style={{ backgroundColor: 'rgba(248,81,73,0.15)' }}>
                <Trash2 size={20} style={{ color: 'var(--destructive)' }} />
              </div>
              <div>
                <h3 className="text-lg font-semibold" style={{ color: 'var(--text-primary)' }}>Clean Teams</h3>
                <p className="text-sm" style={{ color: 'var(--text-muted)' }}>Remove stale or all teams from ~/.claude/teams/</p>
              </div>
            </div>
            <div className="space-y-4">
              <label className="flex items-center gap-2 text-sm" style={{ color: 'var(--text-secondary)' }}>
                <input
                  type="checkbox"
                  checked={teamsAll}
                  onChange={(e) => setTeamsAll(e.target.checked)}
                />
                Remove all teams (not just stale ones)
              </label>
              <label className="flex items-center gap-2 text-sm" style={{ color: 'var(--text-secondary)' }}>
                <input
                  type="checkbox"
                  checked={teamsDryRun}
                  onChange={(e) => setTeamsDryRun(e.target.checked)}
                />
                Dry run (preview only)
              </label>
              {teamsDryRun && (
                <div className="flex items-center gap-2 text-sm" style={{ color: 'var(--warning)' }}>
                  <AlertTriangle size={14} />
                  No teams will be deleted in dry run mode
                </div>
              )}
              <button
                onClick={cleanTeams}
                disabled={loading}
                className="btn-primary flex items-center gap-2 whitespace-nowrap"
                style={{ backgroundColor: 'var(--destructive)' }}
              >
                {loading ? (
                  <RefreshCw size={16} className="animate-spin" />
                ) : (
                  <Trash2 size={16} />
                )}
                {teamsDryRun ? 'Preview Clean' : 'Clean Teams'}
              </button>
            </div>
          </div>

          {teamsResult && (
            <div className="card p-4">
              <div className="flex items-center gap-2 mb-3">
                {teamsResult.dry_run ? (
                  <>
                    <AlertTriangle size={16} style={{ color: 'var(--warning)' }} />
                    <h4 className="text-sm font-semibold" style={{ color: 'var(--warning)' }}>Dry Run Preview</h4>
                  </>
                ) : (
                  <>
                    <CheckCircle2 size={16} style={{ color: 'var(--success)' }} />
                    <h4 className="text-sm font-semibold" style={{ color: 'var(--success)' }}>Clean Complete</h4>
                  </>
                )}
              </div>
              <p className="text-sm mb-2" style={{ color: 'var(--text-primary)' }}>{teamsResult.message}</p>
              <p className="text-sm" style={{ color: 'var(--text-muted)' }}>
                {teamsResult.teams_count} teams would be/were removed
              </p>
              {teamsResult.team_names && teamsResult.team_names.length > 0 && (
                <div className="mt-3 flex flex-wrap gap-2">
                  {teamsResult.team_names.map((name) => (
                    <span key={name} className="badge badge-destructive">{name}</span>
                  ))}
                </div>
              )}
            </div>
          )}
        </>
      )}

      {error && (
        <div className="card p-4 border-[var(--destructive)]">
          <p className="text-sm" style={{ color: 'var(--destructive)' }}>{error}</p>
        </div>
      )}
    </div>
  );
}
