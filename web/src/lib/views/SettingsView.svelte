<script lang="ts">
  import { onMount } from 'svelte';
  import { api, type ConfigData } from '../api-client';
  import { Settings, Save, RefreshCw, CheckCircle2, FolderOpen, Plus, Trash2 } from 'lucide-svelte';
  import FileSystemBrowser from '../components/FileSystemBrowser.svelte';
  import PathInput from '../components/PathInput.svelte';

  let config = $state<ConfigData | null>(null);
  let loading = $state(true);
  let saving = $state(false);
  let error = $state<string | null>(null);
  let success = $state(false);

  // Form fields
  let scanPaths = $state<string[]>([]);
  let excludePaths = $state('');
  let maxDepth = $state(5);
  let defaultOutput = $state('table');

  // File browser
  let showBrowser = $state(false);
  let browserTargetPath = $state<'scan' | 'exclude'>('scan');

  async function loadConfig() {
    loading = true;
    error = null;
    try {
      const data = await api.getConfig();
      config = data.config;
      scanPaths = config?.scan_paths ?? [];
      excludePaths = config?.exclude_paths?.join('\n') ?? '';
      maxDepth = config?.max_depth ?? 5;
      defaultOutput = config?.default_output ?? 'table';
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  }

  async function saveConfig() {
    saving = true;
    error = null;
    success = false;
    try {
      await api.saveConfig({
        scan_paths: scanPaths,
        exclude_paths: excludePaths.split('\n').map(s => s.trim()).filter(Boolean),
        max_depth: maxDepth,
        default_output: defaultOutput,
      });
      success = true;
      setTimeout(() => success = false, 3000);
    } catch (e: any) {
      error = e.message;
    } finally {
      saving = false;
    }
  }

  function openBrowser(target: 'scan' | 'exclude' = 'scan') {
    browserTargetPath = target;
    showBrowser = true;
  }

  function handlePathSelect(event: CustomEvent<{ path: string }>) {
    const selectedPath = event.detail.path;
    if (browserTargetPath === 'scan') {
      if (!scanPaths.includes(selectedPath)) {
        scanPaths = [...scanPaths, selectedPath];
      }
    } else if (browserTargetPath === 'exclude') {
      excludePaths = excludePaths ? `${excludePaths}\n${selectedPath}` : selectedPath;
    }
    showBrowser = false;
  }

  function removeScanPath(pathToRemove: string) {
    scanPaths = scanPaths.filter(p => p !== pathToRemove);
  }

  function removeExcludePath(index: number) {
    const lines = excludePaths.split('\n');
    lines.splice(index, 1);
    excludePaths = lines.join('\n');
  }

  onMount(() => {
    loadConfig();
  });
</script>

<div class="space-y-6">
  <div class="card p-6">
    <div class="flex items-center gap-3 mb-6">
      <div class="w-12 h-12 rounded-xl flex items-center justify-center" style="background-color: var(--accent);">
        <Settings size={24} class="text-white" />
      </div>
      <div>
        <h2 class="text-xl font-semibold" style="color: var(--text-primary);">Settings</h2>
        <p class="text-sm" style="color: var(--text-muted);">Configure ccinspect scanning behavior</p>
      </div>
    </div>

    {#if loading}
      <div class="text-center py-12" style="color: var(--text-muted);">
        <RefreshCw size={32} class="animate-spin mx-auto mb-4" />
        <p>Loading config...</p>
      </div>
    {:else if error && !config}
      <div class="card p-4 border-[var(--destructive)]">
        <p class="text-sm" style="color: var(--destructive);">Error: {error}</p>
        <button onclick={loadConfig} class="btn-secondary mt-2 text-sm">Retry</button>
      </div>
    {:else}
      <div class="space-y-6">
        <!-- Scan Paths -->
        <div>
          <label class="block text-sm font-medium mb-2" style="color: var(--text-primary);">Scan Paths</label>
          <div class="space-y-2">
            {#each scanPaths as path, i}
              <div class="flex items-center gap-2">
                <div class="flex-1 flex items-center gap-2 px-3 py-2 rounded-lg font-mono text-sm"
                     style="background-color: var(--bg-muted); color: var(--text-secondary); border: 1px solid var(--border-subtle);">
                  <FolderOpen size={14} style="color: var(--accent);" />
                  <span class="truncate">{path}</span>
                </div>
                <button class="btn-icon" onclick={() => removeScanPath(path)} title="Remove">
                  <span style="font-size: 16px;">×</span>
                </button>
              </div>
            {/each}
            <button onclick={openBrowser} class="btn-secondary flex items-center gap-2 w-full">
              <Plus size={16} />
              Add Folder...
            </button>
          </div>
          <p class="text-xs mt-1" style="color: var(--text-muted);">Select directories to scan for .claude/ directories.</p>
        </div>

        <!-- Exclude Paths -->
        <div>
          <label class="block text-sm font-medium mb-2" style="color: var(--text-primary);">Exclude Paths</label>
          <div class="space-y-2">
            {#each excludePaths.split('\n').filter(Boolean) as line, i}
              <div class="flex items-center gap-2">
                <div class="flex-1 flex items-center gap-2 px-3 py-2 rounded-lg font-mono text-sm"
                     style="background-color: var(--bg-muted); color: var(--text-secondary); border: 1px solid var(--border-subtle);">
                  <Trash2 size={14} style="color: var(--destructive);" />
                  <span class="truncate">{line}</span>
                </div>
                <button class="btn-icon" onclick={() => removeExcludePath(i)} title="Remove">
                  <span style="font-size: 16px;">×</span>
                </button>
              </div>
            {/each}
            <button onclick={() => openBrowser('exclude')} class="btn-secondary flex items-center gap-2 w-full">
              <Plus size={16} />
              Add Exclude Path...
            </button>
          </div>
          <p class="text-xs mt-1" style="color: var(--text-muted);">Directory names to skip during scanning.</p>
        </div>

        <!-- Max Depth -->
        <div>
          <label class="block text-sm font-medium mb-2" style="color: var(--text-primary);">Max Depth</label>
          <input
            type="number"
            bind:value={maxDepth}
            min="1"
            max="20"
            class="input-field w-32"
          />
          <p class="text-xs mt-1" style="color: var(--text-muted);">How deep to recurse into scan paths.</p>
        </div>

        <!-- Default Output -->
        <div>
          <label class="block text-sm font-medium mb-2" style="color: var(--text-primary);">Default Output</label>
          <select bind:value={defaultOutput} class="input-field input-field-select w-40">
            <option value="table">Table</option>
            <option value="json">JSON</option>
            <option value="markdown">Markdown</option>
          </select>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-4">
          <button onclick={saveConfig} disabled={saving} class="btn-primary flex items-center gap-2">
            {#if saving}
              <RefreshCw size={16} class="animate-spin" />
              Saving...
            {:else}
              <Save size={16} />
              Save Settings
            {/if}
          </button>
          {#if success}
            <span class="flex items-center gap-1 text-sm" style="color: var(--success);">
              <CheckCircle2 size={14} /> Saved!
            </span>
          {/if}
        </div>

        {#if error}
          <div class="card p-4 border-[var(--destructive)]">
            <p class="text-sm" style="color: var(--destructive);">{error}</p>
          </div>
        {/if}
      </div>
    {/if}
  </div>

  <!-- File Browser Modal -->
  {#if showBrowser}
    <FileSystemBrowser
      selectFolder={true}
      onselect={handlePathSelect}
      onclose={() => showBrowser = false}
    />
  {/if}
</div>
