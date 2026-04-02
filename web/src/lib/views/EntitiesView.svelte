<script lang="ts">
  import { onMount } from 'svelte';
  import { api, type ListResponse } from '../api-client';
  import {
    Package, Webhook, Bot, Command, FileText, Server, Users, Copy, ChevronRight, RefreshCw, Filter
  } from 'lucide-svelte';

  let entityType = $state('skills');
  let prefix = $state('');
  let globalOnly = $state(false);
  let listData = $state<ListResponse | null>(null);
  let loading = $state(false);
  let error = $state<string | null>(null);

  const entityFilters = ['skills', 'hooks', 'agents', 'commands', 'rules', 'mcp', 'teams'];

  const entityIcons: Record<string, any> = {
    skill: Package,
    hook: Webhook,
    agent: Bot,
    command: Command,
    rule: FileText,
    mcp: Server,
    team: Users,
  };

  async function loadEntities() {
    loading = true;
    error = null;
    try {
      listData = await api.list(entityType, { global: globalOnly, prefix });
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  function getIcon(type: string) {
    return entityIcons[type] || Package;
  }

  function shortenPath(path: string): string {
    if (path === 'global') return 'global';
    const parts = path.split('/');
    return '~/' + parts.slice(-2).join('/');
  }

  onMount(() => {
    loadEntities();
  });
</script>

<div class="space-y-6">
  <!-- Filter Bar -->
  <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
    <div class="flex items-center gap-2 overflow-x-auto pb-2">
      {#each entityFilters as filter}
        <button
          onclick={() => { entityType = filter; loadEntities(); }}
          class="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-all duration-200"
          style="background-color: {entityType === filter ? 'var(--accent)' : 'var(--bg-muted)'}; color: {entityType === filter ? 'white' : 'var(--text-secondary)'};"
        >
          {filter.charAt(0).toUpperCase() + filter.slice(1)}
        </button>
      {/each}
    </div>
    <div class="flex flex-wrap items-center gap-2">
      <input
        type="text"
        bind:value={prefix}
        placeholder="Prefix filter..."
        class="input-field w-40 text-sm"
        onkeydown={(e) => e.key === 'Enter' && loadEntities()}
      />
      <label class="flex items-center gap-2 text-sm whitespace-nowrap" style="color: var(--text-secondary);">
        <input type="checkbox" bind:checked={globalOnly} onchange={loadEntities} />
        Global only
      </label>
      <button onclick={loadEntities} disabled={loading} class="btn-secondary text-sm">
        <RefreshCw size={14} class="{loading ? 'animate-spin' : ''}" />
      </button>
    </div>
  </div>

  {#if error}
    <div class="card p-4 border-[var(--destructive)]">
      <p class="text-sm" style="color: var(--destructive);">Error: {error}</p>
    </div>
  {:else if loading && !listData}
    <div class="text-center py-12" style="color: var(--text-muted);">
      <RefreshCw size={32} class="animate-spin mx-auto mb-4" />
      <p>Loading entities...</p>
    </div>
  {:else}
    <!-- Entities List -->
    <div class="space-y-2">
      {#each listData?.entities ?? [] as entity, i}
        {@const Icon = getIcon(entity.type ?? entityType)}
        <div class="entity-list-item animate-in" style="animation-delay: {i * 0.03}s;">
          <div class="entity-list-item-icon">
            <Icon size={18} />
          </div>
          <div class="entity-list-item-content">
            <p class="entity-list-item-title">{entity.name}</p>
            {#if entity.description}
              <p class="entity-list-item-meta">{entity.description}</p>
            {/if}
          </div>
          <span class="badge {entity.source === 'global' ? 'badge-brand' : 'badge-success'}">
            {shortenPath(entity.source)}
          </span>
          <div class="flex items-center gap-1">
            <button class="btn-icon w-8 h-8" aria-label="Copy" title="Copy entity">
              <Copy size={14} />
            </button>
            <button class="btn-icon w-8 h-8" aria-label="View details">
              <ChevronRight size={14} />
            </button>
          </div>
        </div>
      {/each}
      {#if listData?.entities.length === 0}
        <div class="text-center py-12" style="color: var(--text-muted);">
          <p>No {entityType} found</p>
        </div>
      {/if}
    </div>
  {/if}
</div>
