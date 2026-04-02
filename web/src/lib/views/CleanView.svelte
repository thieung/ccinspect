<script lang="ts">
  import { api, type CleanResponse, type CleanTeamsResponse } from '../api-client';
  import { Trash2, RefreshCw, AlertTriangle, CheckCircle2 } from 'lucide-svelte';
  import PathInput from '../components/PathInput.svelte';

  let mode: 'project' | 'teams' = $state('project');

  // Project clean
  let projectPath = $state('');
  let projectDryRun = $state(true);
  let projectResult = $state<CleanResponse | null>(null);

  // Teams clean
  let teamsDryRun = $state(true);
  let teamsAll = $state(false);
  let teamsResult = $state<CleanTeamsResponse | null>(null);

  let loading = $state(false);
  let error = $state<string | null>(null);

  async function cleanProject() {
    if (!projectPath) {
      error = 'Please specify a project path';
      return;
    }
    loading = true;
    error = null;
    try {
      projectResult = await api.clean({ path: projectPath, dry_run: projectDryRun });
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  async function cleanTeams() {
    loading = true;
    error = null;
    try {
      teamsResult = await api.cleanTeams({ all: teamsAll, dry_run: teamsDryRun });
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  function handleProjectPathChange(event: CustomEvent<{ path: string }>) {
    projectPath = event.detail.path;
  }
</script>

<div class="space-y-6">
  <!-- Mode Toggle -->
  <div class="flex gap-2">
    <button
      onclick={() => mode = 'project'}
      class="px-4 py-2 rounded-lg text-sm font-medium transition-all"
      style="background-color: {mode === 'project' ? 'var(--accent)' : 'var(--bg-muted)'}; color: {mode === 'project' ? 'white' : 'var(--text-secondary)'};"
    >
      Clean Project
    </button>
    <button
      onclick={() => mode = 'teams'}
      class="px-4 py-2 rounded-lg text-sm font-medium transition-all"
      style="background-color: {mode === 'teams' ? 'var(--accent)' : 'var(--bg-muted)'}; color: {mode === 'teams' ? 'white' : 'var(--text-secondary)'};"
    >
      Clean Teams
    </button>
  </div>

  {#if mode === 'project'}
    <!-- Project Clean -->
    <div class="card p-6">
      <div class="flex items-center gap-3 mb-4">
        <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: rgba(248,81,73,0.15);">
          <Trash2 size={20} style="color: var(--destructive);" />
        </div>
        <div>
          <h3 class="text-lg font-semibold" style="color: var(--text-primary);">Clean Project</h3>
          <p class="text-sm" style="color: var(--text-muted);">Remove .claude/ directory from a project</p>
        </div>
      </div>
      <div class="space-y-4">
        <PathInput value={projectPath} placeholder="/path/to/project" allowGlobal={false} label="Project Path" on:change={handleProjectPathChange} />
        <label class="flex items-center gap-2 text-sm" style="color: var(--text-secondary);">
          <input type="checkbox" bind:checked={projectDryRun} />
          Dry run (preview only)
        </label>
        {#if projectDryRun}
          <div class="flex items-center gap-2 text-sm" style="color: var(--warning);">
            <AlertTriangle size={14} />
            No files will be deleted in dry run mode
          </div>
        {/if}
        <button onclick={cleanProject} disabled={loading} class="btn-primary flex items-center gap-2 whitespace-nowrap" style="background-color: var(--destructive);">
          {#if loading}
            <RefreshCw size={16} class="animate-spin" />
          {:else}
            <Trash2 size={16} />
          {/if}
          {projectDryRun ? 'Preview Clean' : 'Clean Project'}
        </button>
      </div>
    </div>

    {#if projectResult}
      <div class="card p-4">
        <div class="flex items-center gap-2 mb-3">
          {#if projectResult.dry_run}
            <AlertTriangle size={16} style="color: var(--warning);" />
            <h4 class="text-sm font-semibold" style="color: var(--warning);">Dry Run Preview</h4>
          {:else}
            <CheckCircle2 size={16} style="color: var(--success);" />
            <h4 class="text-sm font-semibold" style="color: var(--success);">Clean Complete</h4>
          {/if}
        </div>
        <p class="text-sm mb-2" style="color: var(--text-primary);">{projectResult.message}</p>
        <p class="text-sm" style="color: var(--text-muted);">{projectResult.files_count} files would be/were removed</p>
        {#if projectResult.files && projectResult.files.length > 0}
          <div class="mt-3 max-h-48 overflow-y-auto">
            {#each projectResult.files as file}
              <p class="text-xs font-mono py-0.5" style="color: var(--text-secondary);">{file}</p>
            {/each}
          </div>
        {/if}
      </div>
    {/if}
  {:else}
    <!-- Teams Clean -->
    <div class="card p-6">
      <div class="flex items-center gap-3 mb-4">
        <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: rgba(248,81,73,0.15);">
          <Trash2 size={20} style="color: var(--destructive);" />
        </div>
        <div>
          <h3 class="text-lg font-semibold" style="color: var(--text-primary);">Clean Teams</h3>
          <p class="text-sm" style="color: var(--text-muted);">Remove stale or all teams from ~/.claude/teams/</p>
        </div>
      </div>
      <div class="space-y-4">
        <label class="flex items-center gap-2 text-sm" style="color: var(--text-secondary);">
          <input type="checkbox" bind:checked={teamsAll} />
          Remove all teams (not just stale ones)
        </label>
        <label class="flex items-center gap-2 text-sm" style="color: var(--text-secondary);">
          <input type="checkbox" bind:checked={teamsDryRun} />
          Dry run (preview only)
        </label>
        {#if teamsDryRun}
          <div class="flex items-center gap-2 text-sm" style="color: var(--warning);">
            <AlertTriangle size={14} />
            No teams will be deleted in dry run mode
          </div>
        {/if}
        <button onclick={cleanTeams} disabled={loading} class="btn-primary flex items-center gap-2 whitespace-nowrap" style="background-color: var(--destructive);">
          {#if loading}
            <RefreshCw size={16} class="animate-spin" />
          {:else}
            <Trash2 size={16} />
          {/if}
          {teamsDryRun ? 'Preview Clean' : 'Clean Teams'}
        </button>
      </div>
    </div>

    {#if teamsResult}
      <div class="card p-4">
        <div class="flex items-center gap-2 mb-3">
          {#if teamsResult.dry_run}
            <AlertTriangle size={16} style="color: var(--warning);" />
            <h4 class="text-sm font-semibold" style="color: var(--warning);">Dry Run Preview</h4>
          {:else}
            <CheckCircle2 size={16} style="color: var(--success);" />
            <h4 class="text-sm font-semibold" style="color: var(--success);">Clean Complete</h4>
          {/if}
        </div>
        <p class="text-sm mb-2" style="color: var(--text-primary);">{teamsResult.message}</p>
        <p class="text-sm" style="color: var(--text-muted);">{teamsResult.teams_count} teams would be/were removed</p>
        {#if teamsResult.team_names && teamsResult.team_names.length > 0}
          <div class="mt-3 flex flex-wrap gap-2">
            {#each teamsResult.team_names as name}
              <span class="badge badge-destructive">{name}</span>
            {/each}
          </div>
        {/if}
      </div>
    {/if}
  {/if}

  {#if error}
    <div class="card p-4 border-[var(--destructive)]">
      <p class="text-sm" style="color: var(--destructive);">{error}</p>
    </div>
  {/if}
</div>
