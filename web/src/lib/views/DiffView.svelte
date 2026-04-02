<script lang="ts">
  import { api, type DiffResponse } from '../api-client';
  import { GitCompare, RefreshCw, Plus, Minus, ArrowRight } from 'lucide-svelte';
  import PathInput from '../components/PathInput.svelte';

  let pathA = $state('global');
  let pathB = $state('');
  let entityType = $state('skills');
  let diffData = $state<DiffResponse | null>(null);
  let loading = $state(false);
  let error = $state<string | null>(null);

  async function runDiff() {
    if (!pathA || !pathB) {
      error = 'Please select two projects to compare';
      return;
    }
    loading = true;
    error = null;
    try {
      diffData = await api.diff(pathA, pathB, entityType);
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  function handlePathAChange(event: CustomEvent<{ path: string }>) {
    pathA = event.detail.path;
  }

  function handlePathBChange(event: CustomEvent<{ path: string }>) {
    pathB = event.detail.path;
  }
</script>

<div class="space-y-6">
  <!-- Diff Controls -->
  <div class="card p-6">
    <h3 class="text-lg font-semibold mb-4" style="color: var(--text-primary);">Compare Projects</h3>
    <div class="flex flex-col md:flex-row gap-4 items-end">
      <div class="flex-1 w-full">
        <PathInput value={pathA} placeholder="global or /path/to/project" on:change={handlePathAChange} />
      </div>
      <div class="flex items-center justify-center py-2">
        <ArrowRight size={20} style="color: var(--text-muted);" />
      </div>
      <div class="flex-1 w-full">
        <PathInput value={pathB} placeholder="/path/to/project" allowGlobal={false} on:change={handlePathBChange} />
      </div>
      <div class="w-32">
        <label class="block text-sm mb-1" style="color: var(--text-secondary);">Type</label>
        <select bind:value={entityType} class="input-field input-field-select">
          <option value="skills">Skills</option>
          <option value="hooks">Hooks</option>
          <option value="agents">Agents</option>
          <option value="commands">Commands</option>
        </select>
      </div>
      <button onclick={runDiff} disabled={loading} class="btn-primary flex items-center gap-2 whitespace-nowrap">
        {#if loading}
          <RefreshCw size={16} class="animate-spin" />
        {:else}
          <GitCompare size={16} />
        {/if}
        Compare
      </button>
    </div>
  </div>

  {#if error}
    <div class="card p-4 border-[var(--destructive)]">
      <p class="text-sm" style="color: var(--destructive);">{error}</p>
    </div>
  {/if}

  {#if diffData}
    <!-- Results Header -->
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-semibold" style="color: var(--text-primary);">
        Diff: {diffData.left_name} vs {diffData.right_name}
      </h3>
      <div class="flex gap-4 text-sm">
        <span class="flex items-center gap-1" style="color: var(--success);">
          <Plus size={14} /> {diffData.added.length} added
        </span>
        <span class="flex items-center gap-1" style="color: var(--destructive);">
          <Minus size={14} /> {diffData.removed.length} removed
        </span>
      </div>
    </div>

    <!-- Added (in B but not in A) -->
    {#if diffData.added.length > 0}
      <div class="card p-4">
        <h4 class="text-sm font-semibold mb-3" style="color: var(--success);">Added in {diffData.right_name}</h4>
        <div class="space-y-1">
          {#each diffData.added as item}
            <div class="flex items-center gap-2 py-1">
              <Plus size={14} style="color: var(--success);" />
              <span class="font-mono text-sm" style="color: var(--text-primary);">{item.name}</span>
              {#if item.description}
                <span class="text-xs" style="color: var(--text-muted);">— {item.description}</span>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Removed (in A but not in B) -->
    {#if diffData.removed.length > 0}
      <div class="card p-4">
        <h4 class="text-sm font-semibold mb-3" style="color: var(--destructive);">Removed from {diffData.left_name}</h4>
        <div class="space-y-1">
          {#each diffData.removed as item}
            <div class="flex items-center gap-2 py-1">
              <Minus size={14} style="color: var(--destructive);" />
              <span class="font-mono text-sm" style="color: var(--text-primary);">{item.name}</span>
              {#if item.description}
                <span class="text-xs" style="color: var(--text-muted);">— {item.description}</span>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    {/if}

    {#if diffData.added.length === 0 && diffData.removed.length === 0}
      <div class="card p-8 text-center" style="color: var(--text-muted);">
        <p>No differences found — both projects have the same {entityType}.</p>
      </div>
    {/if}
  {/if}
</div>
