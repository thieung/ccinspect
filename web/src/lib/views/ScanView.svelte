<script lang="ts">
  import { onMount } from 'svelte';
  import { api, type ScanResponse } from '../api-client';
  import {
    Search, Scan, LayoutGrid, List as ListIcon, Terminal, Package, Webhook, Bot, Server, Activity,
    CheckCircle2, RefreshCw, ChevronRight, FolderSearch
  } from 'lucide-svelte';

  interface Props {
    darkMode?: boolean;
  }
  let { darkMode = true }: Props = $props();

  let scanData = $state<ScanResponse | null>(null);
  let loading = $state(false);
  let error = $state<string | null>(null);
  let viewMode: 'grid' | 'list' = $state('grid');
  let searchQuery = $state('');

  async function loadScan() {
    loading = true;
    error = null;
    try {
      scanData = await api.scan();
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  function getProjectName(path: string): string {
    const parts = path.split('/');
    return parts[parts.length - 1] || path;
  }

  let filteredProjects = $derived(
    scanData?.projects.filter(p =>
      getProjectName(p.path).toLowerCase().includes(searchQuery.toLowerCase())
    ) ?? []
  );

  let totalSkills = $derived(scanData?.projects.reduce((s, p) => s + p.skills, 0) ?? 0);
  let totalHooks = $derived(scanData?.projects.reduce((s, p) => s + p.hooks, 0) ?? 0);
  let totalAgents = $derived(scanData?.projects.reduce((s, p) => s + p.agents, 0) ?? 0);
  let totalMcp = $derived(scanData?.projects.reduce((s, p) => s + p.mcp_servers, 0) ?? 0);

  onMount(() => {
    loadScan();
  });
</script>

<div class="space-y-8">
  <!-- Stats Grid -->
  <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
    <div class="stat-card animate-in stagger-1">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
          <FolderSearch size={16} style="color: var(--accent);" />
        </div>
        <span class="text-xs" style="color: var(--text-muted);">Projects</span>
      </div>
      <p class="text-2xl font-bold" style="color: var(--text-primary);">{scanData?.projects.length ?? 0}</p>
    </div>
    <div class="stat-card animate-in stagger-2">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
          <Package size={16} style="color: var(--accent);" />
        </div>
        <span class="text-xs" style="color: var(--text-muted);">Skills</span>
      </div>
      <p class="text-2xl font-bold" style="color: var(--text-primary);">{totalSkills}</p>
    </div>
    <div class="stat-card animate-in stagger-3">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
          <Webhook size={16} style="color: var(--accent);" />
        </div>
        <span class="text-xs" style="color: var(--text-muted);">Hooks</span>
      </div>
      <p class="text-2xl font-bold" style="color: var(--text-primary);">{totalHooks}</p>
    </div>
    <div class="stat-card animate-in stagger-4">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
          <Bot size={16} style="color: var(--accent);" />
        </div>
        <span class="text-xs" style="color: var(--text-muted);">Agents</span>
      </div>
      <p class="text-2xl font-bold" style="color: var(--text-primary);">{totalAgents}</p>
    </div>
    <div class="stat-card animate-in stagger-5">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
          <Server size={16} style="color: var(--accent);" />
        </div>
        <span class="text-xs" style="color: var(--text-muted);">MCP</span>
      </div>
      <p class="text-2xl font-bold" style="color: var(--text-primary);">{totalMcp}</p>
    </div>
    <div class="stat-card animate-in stagger-6">
      <div class="flex items-center gap-2 mb-2">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
          <Activity size={16} style="color: var(--accent);" />
        </div>
        <span class="text-xs" style="color: var(--text-muted);">Status</span>
      </div>
      <p class="text-lg font-bold" style="color: var(--success);">{loading ? 'Scanning' : 'Active'}</p>
    </div>
  </div>

  {#if error}
    <div class="card p-4 border-[var(--destructive)]">
      <p class="text-sm" style="color: var(--destructive);">Error: {error}</p>
      <button onclick={loadScan} class="btn-secondary mt-2 text-sm">Retry</button>
    </div>
  {/if}

  <!-- Action Bar -->
  <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
    <div class="relative w-full sm:w-80">
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Search projects..."
        class="input-field pl-11"
      />
      <Search size={18} class="absolute left-3 top-1/2 -translate-y-1/2 z-10" style="color: var(--text-muted);" />
    </div>
    <div class="flex items-center gap-2">
      <button onclick={loadScan} disabled={loading} class="btn-primary flex items-center gap-2 whitespace-nowrap">
        {#if loading}
          <RefreshCw size={16} class="animate-spin" />
          Scanning...
        {:else}
          <Scan size={16} />
          Scan Projects
        {/if}
      </button>
      <div class="flex items-center gap-1 p-1 rounded-lg" style="background-color: var(--bg-muted);">
        <button
          onclick={() => viewMode = 'grid'}
          class="p-2 rounded-md transition-all duration-200"
          style="background-color: {viewMode === 'grid' ? 'var(--bg-card)' : 'transparent'}; color: {viewMode === 'grid' ? 'var(--accent)' : 'var(--text-muted)'};"
          aria-label="Grid view"
        >
          <LayoutGrid size={18} />
        </button>
        <button
          onclick={() => viewMode = 'list'}
          class="p-2 rounded-md transition-all duration-200"
          style="background-color: {viewMode === 'list' ? 'var(--bg-card)' : 'transparent'}; color: {viewMode === 'list' ? 'var(--accent)' : 'var(--text-muted)'};"
          aria-label="List view"
        >
          <ListIcon size={18} />
        </button>
      </div>
    </div>
  </div>

  <!-- Projects Grid/List -->
  {#if loading && !scanData}
    <div class="text-center py-12" style="color: var(--text-muted);">
      <RefreshCw size={32} class="animate-spin mx-auto mb-4" />
      <p>Scanning projects...</p>
    </div>
  {:else if viewMode === 'grid'}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {#each filteredProjects as project, i}
        <div class="card p-5 cursor-pointer animate-in" style="animation-delay: {i * 0.05}s;">
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent);">
                <Terminal size={20} class="text-white" />
              </div>
              <div>
                <h3 class="font-semibold" style="color: var(--text-primary);">{getProjectName(project.path)}</h3>
                <p class="text-xs font-mono truncate max-w-48" style="color: var(--text-muted);">{project.path}</p>
              </div>
            </div>
            <button class="btn-icon w-8 h-8" aria-label="Project menu">
              <ChevronRight size={16} />
            </button>
          </div>
          <div class="grid grid-cols-4 gap-2 mb-3">
            <div class="text-center p-2 rounded-lg" style="background-color: var(--bg-muted);">
              <p class="text-lg font-bold" style="color: var(--accent);">{project.skills}</p>
              <p class="text-xs" style="color: var(--text-muted);">Skills</p>
            </div>
            <div class="text-center p-2 rounded-lg" style="background-color: var(--bg-muted);">
              <p class="text-lg font-bold" style="color: var(--accent);">{project.hooks}</p>
              <p class="text-xs" style="color: var(--text-muted);">Hooks</p>
            </div>
            <div class="text-center p-2 rounded-lg" style="background-color: var(--bg-muted);">
              <p class="text-lg font-bold" style="color: var(--accent);">{project.agents}</p>
              <p class="text-xs" style="color: var(--text-muted);">Agents</p>
            </div>
            <div class="text-center p-2 rounded-lg" style="background-color: var(--bg-muted);">
              <p class="text-lg font-bold" style="color: var(--accent);">{project.mcp_servers}</p>
              <p class="text-xs" style="color: var(--text-muted);">MCP</p>
            </div>
          </div>
          <div class="flex items-center justify-between">
            <span class="badge {project.has_claude_md ? 'badge-success' : 'badge-warning'}">
              <CheckCircle2 size={12} class="mr-1" />
              {project.has_claude_md ? 'Has CLAUDE.md' : 'No CLAUDE.md'}
            </span>
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <div class="space-y-2">
      {#each filteredProjects as project, i}
        <div class="entity-list-item animate-in" style="animation-delay: {i * 0.03}s;">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent);">
            <Terminal size={18} class="text-white" />
          </div>
          <div class="entity-list-item-content">
            <p class="entity-list-item-title">{getProjectName(project.path)}</p>
            <p class="entity-list-item-meta font-mono text-xs">{project.path}</p>
          </div>
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-3 text-sm">
              <span class="flex items-center gap-1" style="color: var(--text-muted);">
                <Package size={14} /> {project.skills}
              </span>
              <span class="flex items-center gap-1" style="color: var(--text-muted);">
                <Webhook size={14} /> {project.hooks}
              </span>
              <span class="flex items-center gap-1" style="color: var(--text-muted);">
                <Bot size={14} /> {project.agents}
              </span>
            </div>
            <span class="badge {project.has_claude_md ? 'badge-success' : 'badge-warning'}">
              {project.has_claude_md ? 'CLAUDE.md' : 'No CLAUDE.md'}
            </span>
          </div>
          <ChevronRight size={18} style="color: var(--text-muted);" />
        </div>
      {/each}
    </div>
  {/if}
</div>
