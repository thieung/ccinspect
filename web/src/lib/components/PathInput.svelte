<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { FolderOpen } from 'lucide-svelte';
  import FileSystemBrowser from './FileSystemBrowser.svelte';

  interface Props {
    value?: string;
    placeholder?: string;
    label?: string;
    allowGlobal?: boolean;
  }

  let { value = '', placeholder = '/path/to/project', label, allowGlobal = true }: Props = $props();
  const dispatch = createEventDispatcher<{ change: string }>();

  let showBrowser = $state(false);
  let inputValue = $state(value);

  function openBrowser() {
    showBrowser = true;
  }

  function handlePathSelect(event: CustomEvent<{ path: string }>) {
    inputValue = event.detail.path;
    dispatch('change', { path: inputValue });
    showBrowser = false;
  }

  function handleInput() {
    dispatch('change', { path: inputValue });
  }
</script>

<div class="space-y-1">
  {#if label}
    <label class="block text-sm font-medium" style="color: var(--text-secondary);">{label}</label>
  {/if}
  <div class="flex items-center gap-2">
    <input
      type="text"
      bind:value={inputValue}
      oninput={handleInput}
      placeholder={placeholder}
      class="input-field flex-1"
    />
    {#if allowGlobal}
      <button
        type="button"
        onclick={() => { inputValue = 'global'; handleInput(); }}
        class="px-3 py-2 text-xs font-medium rounded-md transition-all"
        style="background-color: var(--bg-muted); color: var(--text-secondary); border: 1px solid var(--border-subtle);"
        title="Use global config"
      >
        global
      </button>
    {/if}
    <button
      type="button"
      onclick={openBrowser}
      class="btn-icon"
      title="Browse folder"
    >
      <FolderOpen size={16} />
    </button>
  </div>
</div>

{#if showBrowser}
  <FileSystemBrowser
    selectFolder={true}
    onselect={handlePathSelect}
    onclose={() => showBrowser = false}
  />
{/if}
