<script lang="ts">
  import { onMount } from 'svelte';
  import { api, type ScanResponse } from '../api-client';
  import { Globe, Package, Webhook, Bot, Server, Copy, GitCompare, Trash2, RefreshCw, FileText, Users } from 'lucide-svelte';

  interface Props {
    onNavigate?: (tab: string) => void;
  }
  let { onNavigate }: Props = $props();

  let scanData = $state<ScanResponse | null>(null);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let expandedSection = $state<'skills' | 'hooks' | 'agents' | 'commands' | 'rules' | 'teams' | null>(null);

  async function loadGlobal() {
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

  onMount(() => {
    loadGlobal();
  });

  function navigate(tab: string) {
    if (onNavigate) onNavigate(tab);
  }

  function toggleSection(section: string) {
    expandedSection = expandedSection === section ? null : (section as any);
  }
</script>

<div class="space-y-6">
  {#if loading}
    <div class="text-center py-12" style="color: var(--text-muted);">
      <RefreshCw size={32} class="animate-spin mx-auto mb-4" />
      <p>Loading global config...</p>
    </div>
  {:else if error}
    <div class="card p-4 border-[var(--destructive)]">
      <p class="text-sm" style="color: var(--destructive);">Error: {error}</p>
      <button onclick={loadGlobal} class="btn-secondary mt-2 text-sm">Retry</button>
    </div>
  {:else}
    <div class="card p-6">
      <div class="flex items-center gap-3 mb-6">
        <div class="w-12 h-12 rounded-xl flex items-center justify-center" style="background-color: var(--accent);">
          <Globe size={24} class="text-white" />
        </div>
        <div>
          <h2 class="text-xl font-semibold" style="color: var(--text-primary);">Global Configuration</h2>
          <p class="text-sm" style="color: var(--text-muted);">{scanData?.global?.path ?? '~/.claude/'}</p>
        </div>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
             style="background-color: var(--bg-muted); border-color: var(--border-color);"
             onclick={() => toggleSection('skills')}>
          <Package size={20} style="color: var(--accent);" class="mb-2" />
          <p class="text-2xl font-bold" style="color: var(--text-primary);">{scanData?.global?.skills ?? 0}</p>
          <p class="text-sm" style="color: var(--text-muted);">Skills</p>
        </div>
        <div class="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
             style="background-color: var(--bg-muted); border-color: var(--border-color);"
             onclick={() => toggleSection('hooks')}>
          <Webhook size={20} style="color: var(--accent);" class="mb-2" />
          <p class="text-2xl font-bold" style="color: var(--text-primary);">{scanData?.global?.hooks ?? 0}</p>
          <p class="text-sm" style="color: var(--text-muted);">Hooks</p>
        </div>
        <div class="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
             style="background-color: var(--bg-muted); border-color: var(--border-color);"
             onclick={() => toggleSection('agents')}>
          <Bot size={20} style="color: var(--accent);" class="mb-2" />
          <p class="text-2xl font-bold" style="color: var(--text-primary);">{scanData?.global?.agents ?? 0}</p>
          <p class="text-sm" style="color: var(--text-muted);">Agents</p>
        </div>
        <div class="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
             style="background-color: var(--bg-muted); border-color: var(--border-color);"
             onclick={() => toggleSection('commands')}>
          <Server size={20} style="color: var(--accent);" class="mb-2" />
          <p class="text-2xl font-bold" style="color: var(--text-primary);">{scanData?.global?.commands ?? 0}</p>
          <p class="text-sm" style="color: var(--text-muted);">Commands</p>
        </div>
        <div class="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
             style="background-color: var(--bg-muted); border-color: var(--border-color);"
             onclick={() => toggleSection('rules')}>
          <FileText size={20} style="color: var(--accent);" class="mb-2" />
          <p class="text-2xl font-bold" style="color: var(--text-primary);">{scanData?.global?.rules ?? 0}</p>
          <p class="text-sm" style="color: var(--text-muted);">Rules</p>
        </div>
        <div class="p-4 rounded-xl border cursor-pointer transition-all hover:border-[var(--accent)]"
             style="background-color: var(--bg-muted); border-color: var(--border-color);"
             onclick={() => toggleSection('teams')}>
          <Users size={20} style="color: var(--accent);" class="mb-2" />
          <p class="text-2xl font-bold" style="color: var(--text-primary);">{scanData?.global?.teams ?? 0}</p>
          <p class="text-sm" style="color: var(--text-muted);">Teams</p>
        </div>
      </div>

      <!-- Expanded Sections -->
      {#if expandedSection}
        <div class="mt-6 p-4 rounded-xl border" style="background-color: var(--bg-card); border-color: var(--border-color);">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-semibold" style="color: var(--text-primary);">
              {expandedSection.charAt(0).toUpperCase() + expandedSection.slice(1)}
            </h3>
            <button class="btn-icon w-6 h-6" onclick={() => expandedSection = null}>
              <span style="font-size: 16px;">×</span>
            </button>
          </div>
          <p class="text-sm" style="color: var(--text-muted);">
            Use the <strong>Entities</strong> tab to view and manage {expandedSection} in detail.
          </p>
          <button onclick={() => navigate('entities')} class="btn-primary mt-3 text-sm">
            View All {expandedSection.charAt(0).toUpperCase() + expandedSection.slice(1)}
          </button>
        </div>
      {/if}
    </div>

    <!-- Quick Actions -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <button onclick={() => navigate('copy')} class="card p-5 flex items-center gap-4 hover:border-[var(--accent)] transition-all duration-200 cursor-pointer">
        <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent-muted);">
          <Copy size={20} style="color: var(--accent);" />
        </div>
        <div class="text-left">
          <p class="font-medium" style="color: var(--text-primary);">Copy Entity</p>
          <p class="text-sm" style="color: var(--text-muted);">Copy to project</p>
        </div>
      </button>
      <button onclick={() => navigate('diff')} class="card p-5 flex items-center gap-4 hover:border-[var(--accent)] transition-all duration-200 cursor-pointer">
        <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent-muted);">
          <GitCompare size={20} style="color: var(--accent);" />
        </div>
        <div class="text-left">
          <p class="font-medium" style="color: var(--text-primary);">Diff Skills</p>
          <p class="text-sm" style="color: var(--text-muted);">Compare projects</p>
        </div>
      </button>
      <button onclick={() => navigate('clean')} class="card p-5 flex items-center gap-4 hover:border-[var(--accent)] transition-all duration-200 cursor-pointer">
        <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent-muted);">
          <Trash2 size={20} style="color: var(--destructive);" />
        </div>
        <div class="text-left">
          <p class="font-medium" style="color: var(--text-primary);">Clean Teams</p>
          <p class="text-sm" style="color: var(--text-muted);">Remove stale teams</p>
        </div>
      </button>
    </div>
  {/if}
</div>
