const API_BASE = '/api';

async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `HTTP ${res.status}`);
  }
  return res.json();
}

async function post<T>(path: string, body: object): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || `HTTP ${res.status}`);
  }
  return res.json();
}

// Types matching Go types
export interface ProjectSummary {
  path: string;
  skills: number;
  hooks: number;
  agents: number;
  commands: number;
  mcp_servers: number;
  teams: number;
  has_claude_md: boolean;
}

export interface GlobalSummary {
  path: string;
  skills: number;
  hooks: number;
  agents: number;
  commands: number;
  rules: number;
  teams: number;
  has_claude_md: boolean;
}

export interface ScanResponse {
  projects: ProjectSummary[];
  global?: GlobalSummary;
}

export interface EntityItem {
  name: string;
  description?: string;
  source: string;
  type?: string;
}

export interface ListResponse {
  entities: EntityItem[];
}

export interface DiffEntry {
  name: string;
  description?: string;
  side?: string;
}

export interface DiffResponse {
  left_name: string;
  right_name: string;
  entity_type: string;
  added: DiffEntry[];
  removed: DiffEntry[];
  changed: DiffEntry[];
}

export interface CopyRequest {
  type: string;
  names: string[];
  from: string;
  to: string;
  force: boolean;
  dry_run: boolean;
}

export interface CopyResult {
  name: string;
  from: string;
  to: string;
  status: string;
  detail?: string;
}

export interface CopyResponse {
  status: string;
  results?: CopyResult[];
  message?: string;
}

export interface CleanRequest {
  path: string;
  dry_run: boolean;
}

export interface CleanResponse {
  status: string;
  message?: string;
  files_count: number;
  files?: string[];
  dry_run: boolean;
}

export interface CleanTeamsRequest {
  team_names?: string[];
  all: boolean;
  dry_run: boolean;
}

export interface CleanTeamsResponse {
  status: string;
  message?: string;
  teams_count: number;
  team_names?: string[];
  dry_run: boolean;
}

export interface ConfigData {
  scan_paths: string[];
  exclude_paths: string[];
  max_depth: number;
  default_output: string;
}

export interface ConfigResponse {
  config: ConfigData;
}

export interface FSBrowseResponse {
  entries: FSEntry[];
  cwd: string;
}

export interface FSEntry {
  name: string;
  path: string;
  is_dir: boolean;
  size?: number;
  mod_time?: number;
}

export interface HomeDirResponse {
  home: string;
  username: string;
  os: string;
}

// API functions
export const api = {
  scan: () => get<ScanResponse>('/scan'),

  global: () => get<{ global: any }>('/global'),

  list: (type: string, opts: { global?: boolean; prefix?: string } = {}) => {
    const params = new URLSearchParams({ type });
    if (opts.global) params.set('global', 'true');
    if (opts.prefix) params.set('prefix', opts.prefix);
    return get<ListResponse>(`/list?${params}`);
  },

  diff: (a: string, b: string, type = 'skills') => {
    const params = new URLSearchParams({ a, b, type });
    return get<DiffResponse>(`/diff?${params}`);
  },

  copy: (req: CopyRequest) => post<CopyResponse>('/copy', req),

  clean: (req: CleanRequest) => post<CleanResponse>('/clean', req),

  cleanTeams: (req: CleanTeamsRequest) => post<CleanTeamsResponse>('/clean/teams', req),

  getConfig: () => get<ConfigResponse>('/config'),

  saveConfig: (cfg: ConfigData) => post<{ status: string }>('/config', cfg),

  browseFS: (path?: string) => {
    const url = path ? `/fs/browse?path=${encodeURIComponent(path)}` : '/fs/browse';
    return get<FSBrowseResponse>(url);
  },

  getHomeDir: () => get<HomeDirResponse>('/fs/home'),
};
