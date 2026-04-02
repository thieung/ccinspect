<script lang="ts">
  import { api, type CopyResponse } from '../api-client';
  import { Copy as CopyIcon, RefreshCw, ArrowRight, CheckCircle2, XCircle } from 'lucide-svelte';
  import PathInput from '../components/PathInput.svelte';

  let entityType = $state('skills');
  let fromPath = $state('global');
  let toPath = $state('');
  let entityNames = $state('all');
  let force = $state(false);
  let dryRun = $state(true);
  let result = $state<CopyResponse | null>(null);
  let loading = $state(false);
  let error = $state<string | null>(null);

  const entityTypes = ['skills', 'hooks', 'agents', 'commands'];

  async function runCopy() {
    if (!fromPath || !toPath) {
      error = 'Please specify both source and destination';
      return;
    }
    loading = true;
    error = null;
    result = null;
    try {
      const names = entityNames === 'all' ? ['all'] : entityNames.split(',').map(n => n.trim());
      result = await api.copy({
        type: entityType,
        from: fromPath,
        to: toPath,
        names,
        force,
        dry_run: dryRun,
      });
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  function getStatusIcon(status: string) {
    if (status === 'copied' || status === 'would copy') return CheckCircle2;
    return XCircle;
  }

  function handleFromPathChange(event: CustomEvent<{ path: string }>) {
    fromPath = event.detail.path;
  }

  function handleToPathChange(event: CustomEvent<{ path: string }>) {
    toPath = event.detail.path;
  }
</script>

<div class="space-y-6">
  <!-- Copy Controls -->
  <div class="card p-6">
    <h3 class="text-lg font-semibold mb-4" style="color: var(--text-primary);">Copy Entity</h3>
    <div class="space-y-4">
      <div class="flex flex-col md:flex-row gap-4 items-end">
        <div class="w-32">
          <label class="block text-sm mb-1" style="color: var(--text-secondary);">Type</label>
          <select bind:value={entityType} class="input-field input-field-select">
            {#each entityTypes as type}
              <option value={type}>{type}</option>
            {/each}
          </select>
        </div>
        <div class="flex-1 w-full">
          <PathInput value={fromPath} placeholder="global or /path" on:change={handleFromPathChange} />
        </div>
        <div class="flex items-center justify-center py-2">
          <ArrowRight size={20} style="color: var(--text-muted);" />
        </div>
        <div class="flex-1 w-full">
          <PathInput value={toPath} placeholder="/path/to/project" allowGlobal={false} on:change={handleToPathChange} />
        </div>
      </div>
      <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
        <div class="flex-1 w-full">
          <label class="block text-sm mb-1" style="color: var(--text-secondary);">Entity names (comma-separated, or "all")</label>
          <input type="text" bind:value={entityNames} placeholder="all" class="input-field" />
        </div>
        <div class="flex items-center gap-4">
          <label class="flex items-center gap-2 text-sm" style="color: var(--text-secondary);">
            <input type="checkbox" bind:checked={force} />
            Force overwrite
          </label>
          <label class="flex items-center gap-2 text-sm" style="color: var(--text-secondary);">
            <input type="checkbox" bind:checked={dryRun} />
            Dry run
          </label>
        </div>
      </div>
      <button onclick={runCopy} disabled={loading} class="btn-primary flex items-center gap-2 whitespace-nowrap w-full md:w-auto">
        {#if loading}
          <RefreshCw size={16} class="animate-spin" />
        {:else}
          <CopyIcon size={16} />
        {/if}
        {dryRun ? 'Preview Copy' : 'Copy Entity'}
      </button>
    </div>
  </div>

  {#if error}
    <div class="card p-4 border-[var(--destructive)]">
      <p class="text-sm" style="color: var(--destructive);">{error}</p>
    </div>
  {/if}

  {#if result}
    <div class="card p-4">
      <h4 class="text-sm font-semibold mb-3" style="color: var(--text-primary);">Results</h4>
      <div class="space-y-1">
        {#each result.results ?? [] as r}
          {@const Icon = getStatusIcon(r.status)}
          <div class="flex items-center gap-2 py-1">
            <Icon size={14} style="color: {r.status === 'copied' || r.status === 'would copy' ? 'var(--success)' : 'var(--destructive)'};" />
            <span class="font-mono text-sm" style="color: var(--text-primary);">{r.name}</span>
            {#if r.detail}
              <span class="text-xs" style="color: var(--text-muted);">— {r.detail}</span>
            {/if}
          </div>
        {/each}
      </div>
      {#if result.message}
        <p class="mt-3 text-sm" style="color: var(--text-secondary);">{result.message}</p>
      {/if}
    </div>
  {/if}
</div>
