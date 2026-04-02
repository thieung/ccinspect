<script lang="ts">
  import { onMount } from 'svelte';
  import {
    FolderSearch, List, Copy, GitCompare, Trash2, Settings,
    Sun, Moon, Scan, Globe,
  } from 'lucide-svelte';

  import ScanView from './lib/views/ScanView.svelte';
  import GlobalView from './lib/views/GlobalView.svelte';
  import EntitiesView from './lib/views/EntitiesView.svelte';
  import DiffView from './lib/views/DiffView.svelte';
  import CopyView from './lib/views/CopyView.svelte';
  import CleanView from './lib/views/CleanView.svelte';
  import SettingsView from './lib/views/SettingsView.svelte';

  let darkMode = $state(true);
  let activeTab: 'scan' | 'global' | 'entities' | 'diff' | 'copy' | 'clean' | 'settings' = $state('scan');

  const tabs = [
    { id: 'scan', label: 'Scan', icon: FolderSearch },
    { id: 'global', label: 'Global', icon: Globe },
    { id: 'entities', label: 'Entities', icon: List },
    { id: 'diff', label: 'Diff', icon: GitCompare },
    { id: 'copy', label: 'Copy', icon: Copy },
    { id: 'clean', label: 'Clean', icon: Trash2 },
    { id: 'settings', label: 'Settings', icon: Settings },
  ] as const;

  function toggleTheme() {
    darkMode = !darkMode;
    if (darkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }

  function navigate(tab: string) {
    activeTab = tab as typeof activeTab;
  }

  onMount(() => {
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      darkMode = true;
      document.documentElement.classList.add('dark');
    }
  });
</script>

<div class="min-h-screen" style="background-color: var(--bg-primary); color: var(--text-primary);">
  <!-- Header -->
  <header class="sticky top-0 z-50 border-b" style="background-color: var(--bg-primary); border-color: var(--border-color);">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent);">
            <Scan size={22} class="text-white" />
          </div>
          <div>
            <h1 class="text-lg font-semibold" style="color: var(--text-primary);">ccinspect</h1>
            <p class="text-xs" style="color: var(--text-muted);">Claude Code Inspector</p>
          </div>
        </div>

        <!-- Navigation Tabs -->
        <nav class="hidden md:flex items-center gap-1 p-1 rounded-xl" style="background-color: var(--bg-muted);">
          {#each tabs as tab}
            <button
              onclick={() => activeTab = tab.id}
              class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200"
              style="background-color: {activeTab === tab.id ? 'var(--bg-card)' : 'transparent'}; color: {activeTab === tab.id ? 'var(--accent)' : 'var(--text-secondary)'}; box-shadow: {activeTab === tab.id ? 'var(--shadow-card)' : 'none'};"
            >
              <div class="flex items-center gap-2">
                <tab.icon size={16} />
                {tab.label}
              </div>
            </button>
          {/each}
        </nav>

        <!-- Actions -->
        <div class="flex items-center gap-2">
          <button
            onclick={toggleTheme}
            class="btn-icon"
            aria-label="Toggle theme"
          >
            {#if darkMode}
              <Sun size={18} />
            {:else}
              <Moon size={18} />
            {/if}
          </button>
        </div>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if activeTab === 'scan'}
      <ScanView {darkMode} />
    {:else if activeTab === 'global'}
      <GlobalView onNavigate={navigate} />
    {:else if activeTab === 'entities'}
      <EntitiesView />
    {:else if activeTab === 'diff'}
      <DiffView />
    {:else if activeTab === 'copy'}
      <CopyView />
    {:else if activeTab === 'clean'}
      <CleanView />
    {:else if activeTab === 'settings'}
      <SettingsView />
    {/if}
  </main>

  <!-- Footer -->
  <footer class="border-t mt-12 py-6" style="border-color: var(--border-color);">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row items-center justify-between gap-4">
        <p class="text-sm" style="color: var(--text-muted);">
          Built with care by <span class="font-medium" style="color: var(--accent);">@thieunv</span>
        </p>
        <div class="flex items-center gap-4">
          <span class="text-xs font-mono" style="color: var(--text-muted);">v1.0.0</span>
        </div>
      </div>
    </div>
  </footer>
</div>
