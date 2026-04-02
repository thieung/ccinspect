<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { Folder, FolderOpen, ChevronRight, Home, HardDrive } from 'lucide-svelte';
  import { api } from '../api-client';

  interface Props {
    initialPath?: string;
    selectFolder?: boolean;
  }

  let { initialPath = '', selectFolder = true }: Props = $props();
  const dispatch = createEventDispatcher<{
    select: { path: string };
    close: void;
  }>();

  let currentPath = $state('');
  let entries = $state<FSEntry[]>([]);
  let loading = $state(false);
  let error = $state<string | null>(null);
  let homeDir = $state<string>('');

  interface FSEntry {
    name: string;
    path: string;
    is_dir: boolean;
    size?: number;
    mod_time?: number;
  }

  async function loadDirectory(path: string) {
    loading = true;
    error = null;
    try {
      const data = await api.browseFS(path);
      entries = data.entries;
      currentPath = data.cwd;
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  async function loadHomeDir() {
    try {
      const data = await api.getHomeDir();
      homeDir = data.home;
      await loadDirectory(homeDir);
    } catch (e: any) {
      console.error('Failed to load home dir:', e);
      error = 'Failed to load home directory';
    }
  }

  function navigateTo(path: string) {
    loadDirectory(path);
  }

  function handleEntryClick(entry: FSEntry) {
    if (entry.is_dir) {
      navigateTo(entry.path);
    } else {
      // For files, select the parent directory
      dispatch('select', { path: entry.path });
    }
  }

  function formatSize(size?: number): string {
    if (!size) return '';
    if (size < 1024) return `${size} B`;
    if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
    if (size < 1024 * 1024 * 1024) return `${(size / (1024 * 1024)).toFixed(1)} MB`;
    return `${(size / (1024 * 1024 * 1024)).toFixed(1)} GB`;
  }

  onMount(async () => {
    // Always load home dir on mount
    await loadHomeDir();
  });
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center" style="background-color: rgba(0, 0, 0, 0.5);">
  <div class="card w-full max-w-2xl max-h-[80vh] flex flex-col" style="background-color: var(--bg-card);">
    <!-- Header -->
    <div class="flex items-center justify-between p-4 border-b" style="border-color: var(--border-color);">
      <div class="flex items-center gap-2">
        <FolderOpen size={20} style="color: var(--accent);" />
        <h3 class="font-semibold" style="color: var(--text-primary);">Browse Folder</h3>
      </div>
      <button class="btn-icon w-8 h-8" onclick={() => dispatch('close')}>
        <span style="font-size: 20px;">×</span>
      </button>
    </div>

    <!-- Navigation Bar -->
    <div class="flex items-center gap-2 p-4 border-b" style="border-color: var(--border-color);">
      <button class="btn-icon" onclick={loadHomeDir} title="Home">
        <Home size={16} />
      </button>
      {#if homeDir}
        <button class="btn-icon" onclick={() => navigateTo(homeDir)} title="Home Directory">
          <HardDrive size={16} />
        </button>
      {/if}
      <div class="flex-1 flex items-center gap-1 px-3 py-2 rounded-lg font-mono text-sm"
           style="background-color: var(--bg-muted); color: var(--text-secondary);">
        <span class="truncate">{currentPath || 'Loading...'}</span>
      </div>
    </div>

    <!-- Error State -->
    {#if error}
      <div class="p-4 text-center" style="color: var(--destructive);">
        <p>{error}</p>
        <button class="btn-secondary mt-2" onclick={() => loadHomeDir()}>Go to Home</button>
      </div>
    {:else if loading}
      <div class="p-8 text-center" style="color: var(--text-muted);">
        <p>Loading...</p>
      </div>
    {:else}
      <!-- File List -->
      <div class="flex-1 overflow-y-auto p-2">
        {#if entries.length === 0}
          <p class="text-center py-4" style="color: var(--text-muted);">No items in this folder</p>
        {:else}
          <div class="space-y-1">
            {#each entries as entry}
              <button
                class="w-full flex items-center gap-3 p-3 rounded-lg hover:border-[var(--accent)] transition-all text-left"
                style="background-color: {entry.is_dir ? 'var(--bg-muted)' : 'var(--bg-card)'}; border: 1px solid var(--border-subtle);"
                onclick={() => handleEntryClick(entry)}
              >
                {#if entry.is_dir}
                  <Folder size={18} style="color: var(--accent);" />
                {:else}
                  <ChevronRight size={18} style="color: var(--text-muted);" />
                {/if}
                <div class="flex-1 min-w-0">
                  <p class="font-medium truncate" style="color: var(--text-primary);">{entry.name}</p>
                  {#if !entry.is_dir && entry.size}
                    <p class="text-xs" style="color: var(--text-muted);">{formatSize(entry.size)}</p>
                  {/if}
                </div>
                {#if entry.is_dir}
                  <ChevronRight size={16} style="color: var(--text-muted);" />
                {/if}
              </button>
            {/each}
          </div>
        {/if}
      </div>
    {/if}

    <!-- Footer Actions -->
    {#if selectFolder}
      <div class="flex items-center justify-end gap-2 p-4 border-t" style="border-color: var(--border-color);">
        <button class="btn-secondary" onclick={() => dispatch('close')}>Cancel</button>
        <button class="btn-primary" onclick={() => dispatch('select', { path: currentPath })}>
          Select This Folder
        </button>
      </div>
    {/if}
  </div>
</div>
