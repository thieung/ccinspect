<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Search,
    FolderSearch,
    List,
    Copy,
    GitCompare,
    Trash2,
    Settings,
    Sun,
    Moon,
    ChevronRight,
    Terminal,
    Zap,
    Package,
    Webhook,
    Bot,
    Command,
    FileText,
    Server,
    Users,
    Activity,
    AlertCircle,
    CheckCircle2,
    XCircle,
    RefreshCw,
    Plus,
    Filter,
    LayoutGrid,
    List as ListIcon,
    Scan,
    Globe,
  } from 'lucide-svelte';

  // State
  let darkMode = $state(true);
  let viewMode: 'grid' | 'list' = $state('grid');
  let activeTab: 'scan' | 'global' | 'entities' = $state('scan');
  let isScanning = $state(false);
  let scanProgress = $state(0);
  let searchQuery = $state('');
  let selectedEntity: string | null = $state(null);
  let selectedProject: string | null = $state(null);

  // Mock data for demo
  let stats = $state({
    totalProjects: 12,
    totalSkills: 47,
    totalHooks: 8,
    totalAgents: 5,
    totalMcp: 3,
  });

  let projects = $state([
    { name: 'my-app', path: '~/projects/my-app', skills: 12, hooks: 2, agents: 1, mcp: 1, lastScanned: '2 min ago' },
    { name: 'ccinspect', path: '~/projects/solo-builder/ccinspect', skills: 24, hooks: 3, agents: 2, mcp: 2, lastScanned: '5 min ago' },
    { name: 'dashboard', path: '~/projects/dashboard', skills: 8, hooks: 1, agents: 1, mcp: 0, lastScanned: '1 hour ago' },
    { name: 'api-service', path: '~/projects/api-service', skills: 15, hooks: 2, agents: 1, mcp: 3, lastScanned: '3 hours ago' },
  ]);

  let entities = $state([
    { type: 'skill', name: 'code-review', count: 8, location: 'global' },
    { type: 'skill', name: 'planner', count: 5, location: 'project' },
    { type: 'hook', name: 'PreToolUse', count: 3, location: 'global' },
    { type: 'agent', name: 'researcher', count: 2, location: 'global' },
    { type: 'command', name: 'deploy', count: 4, location: 'project' },
    { type: 'mcp', name: 'filesystem', count: 1, location: 'project' },
  ]);

  const entityIcons: Record<string, typeof Zap> = {
    skill: Package,
    hook: Webhook,
    agent: Bot,
    command: Command,
    rule: FileText,
    mcp: Server,
    team: Users,
  };

  // Toggle dark mode
  function toggleTheme() {
    darkMode = !darkMode;
    if (darkMode) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }

  // Scan action
  async function handleScan() {
    isScanning = true;
    scanProgress = 0;
    for (let i = 0; i <= 100; i += 10) {
      await new Promise(r => setTimeout(r, 150));
      scanProgress = i;
    }
    isScanning = false;
  }

  // Entity type filter
  let activeEntityFilter = $state('all');
  const entityFilters = ['all', 'skill', 'hook', 'agent', 'command', 'rule', 'mcp', 'team'];

  let filteredEntities = $derived(
    activeEntityFilter === 'all'
      ? entities
      : entities.filter(e => e.type === activeEntityFilter)
  );

  // Filter projects by search
  let filteredProjects = $derived(
    searchQuery
      ? projects.filter(p => p.name.toLowerCase().includes(searchQuery.toLowerCase()))
      : projects
  );

  onMount(() => {
    // Check system preference
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
          <button
            onclick={() => activeTab = 'scan'}
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200"
            style="background-color: {activeTab === 'scan' ? 'var(--bg-card)' : 'transparent'}; color: {activeTab === 'scan' ? 'var(--accent)' : 'var(--text-secondary)'}; box-shadow: {activeTab === 'scan' ? 'var(--shadow-card)' : 'none'};"
          >
            <div class="flex items-center gap-2">
              <FolderSearch size={16} />
              Scan
            </div>
          </button>
          <button
            onclick={() => activeTab = 'global'}
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200"
            style="background-color: {activeTab === 'global' ? 'var(--bg-card)' : 'transparent'}; color: {activeTab === 'global' ? 'var(--accent)' : 'var(--text-secondary)'}; box-shadow: {activeTab === 'global' ? 'var(--shadow-card)' : 'none'};"
          >
            <div class="flex items-center gap-2">
              <Globe size={16} />
              Global
            </div>
          </button>
          <button
            onclick={() => activeTab = 'entities'}
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200"
            style="background-color: {activeTab === 'entities' ? 'var(--bg-card)' : 'transparent'}; color: {activeTab === 'entities' ? 'var(--accent)' : 'var(--text-secondary)'}; box-shadow: {activeTab === 'entities' ? 'var(--shadow-card)' : 'none'};"
          >
            <div class="flex items-center gap-2">
              <List size={16} />
              Entities
            </div>
          </button>
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
          <button class="btn-icon">
            <Settings size={18} />
          </button>
        </div>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if activeTab === 'scan'}
      <!-- Scan View -->
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
            <p class="text-2xl font-bold" style="color: var(--text-primary);">{stats.totalProjects}</p>
          </div>
          <div class="stat-card animate-in stagger-2">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
                <Package size={16} style="color: var(--accent);" />
              </div>
              <span class="text-xs" style="color: var(--text-muted);">Skills</span>
            </div>
            <p class="text-2xl font-bold" style="color: var(--text-primary);">{stats.totalSkills}</p>
          </div>
          <div class="stat-card animate-in stagger-3">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
                <Webhook size={16} style="color: var(--accent);" />
              </div>
              <span class="text-xs" style="color: var(--text-muted);">Hooks</span>
            </div>
            <p class="text-2xl font-bold" style="color: var(--text-primary);">{stats.totalHooks}</p>
          </div>
          <div class="stat-card animate-in stagger-4">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
                <Bot size={16} style="color: var(--accent);" />
              </div>
              <span class="text-xs" style="color: var(--text-muted);">Agents</span>
            </div>
            <p class="text-2xl font-bold" style="color: var(--text-primary);">{stats.totalAgents}</p>
          </div>
          <div class="stat-card animate-in stagger-5">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
                <Server size={16} style="color: var(--accent);" />
              </div>
              <span class="text-xs" style="color: var(--text-muted);">MCP</span>
            </div>
            <p class="text-2xl font-bold" style="color: var(--text-primary);">{stats.totalMcp}</p>
          </div>
          <div class="stat-card animate-in stagger-6">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background-color: var(--accent-muted);">
                <Activity size={16} style="color: var(--accent);" />
              </div>
              <span class="text-xs" style="color: var(--text-muted);">Status</span>
            </div>
            <p class="text-lg font-bold" style="color: var(--success);">Active</p>
          </div>
        </div>

        <!-- Action Bar -->
        <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
          <div class="relative w-full sm:w-80">
            <Search size={18} class="absolute left-3 top-1/2 -translate-y-1/2" style="color: var(--text-muted);" />
            <input
              type="text"
              bind:value={searchQuery}
              placeholder="Search projects..."
              class="input-field pl-10"
            />
          </div>
          <div class="flex items-center gap-2">
            <button
              onclick={handleScan}
              disabled={isScanning}
              class="btn-primary flex items-center gap-2"
            >
              {#if isScanning}
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

        <!-- Scan Progress -->
        {#if isScanning}
          <div class="card p-4 animate-in">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium" style="color: var(--text-primary);">Scanning directories...</span>
              <span class="text-sm font-mono" style="color: var(--accent);">{scanProgress}%</span>
            </div>
            <div class="h-2 rounded-full overflow-hidden" style="background-color: var(--bg-muted);">
              <div
                class="h-full rounded-full transition-all duration-300 ease-out"
                style="width: {scanProgress}%; background-color: var(--accent);"
              ></div>
            </div>
          </div>
        {/if}

        <!-- Projects Grid/List -->
        {#if viewMode === 'grid'}
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each filteredProjects as project, i}
              <div class="card p-5 cursor-pointer animate-in" style="animation-delay: {i * 0.05}s;">
                <div class="flex items-start justify-between mb-4">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent);">
                      <Terminal size={20} class="text-white" />
                    </div>
                    <div>
                      <h3 class="font-semibold" style="color: var(--text-primary);">{project.name}</h3>
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
                    <p class="text-lg font-bold" style="color: var(--accent);">{project.mcp}</p>
                    <p class="text-xs" style="color: var(--text-muted);">MCP</p>
                  </div>
                </div>
                <div class="flex items-center justify-between">
                  <span class="badge badge-brand">
                    <CheckCircle2 size={12} class="mr-1" />
                    Scanned
                  </span>
                  <span class="text-xs" style="color: var(--text-muted);">{project.lastScanned}</span>
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
                  <p class="entity-list-item-title">{project.name}</p>
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
                  <span class="badge badge-success">Active</span>
                </div>
                <ChevronRight size={18} style="color: var(--text-muted);" />
              </div>
            {/each}
          </div>
        {/if}
      </div>

    {:else if activeTab === 'global'}
      <!-- Global View -->
      <div class="space-y-6">
        <div class="card p-6">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-12 h-12 rounded-xl flex items-center justify-center" style="background-color: var(--accent);">
              <Globe size={24} class="text-white" />
            </div>
            <div>
              <h2 class="text-xl font-semibold" style="color: var(--text-primary);">Global Configuration</h2>
              <p class="text-sm" style="color: var(--text-muted);">~/.claude/</p>
            </div>
          </div>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div class="p-4 rounded-xl border" style="background-color: var(--bg-muted); border-color: var(--border-color);">
              <Package size={20} style="color: var(--accent);" class="mb-2" />
              <p class="text-2xl font-bold" style="color: var(--text-primary);">28</p>
              <p class="text-sm" style="color: var(--text-muted);">Skills</p>
            </div>
            <div class="p-4 rounded-xl border" style="background-color: var(--bg-muted); border-color: var(--border-color);">
              <Webhook size={20} style="color: var(--accent);" class="mb-2" />
              <p class="text-2xl font-bold" style="color: var(--text-primary);">5</p>
              <p class="text-sm" style="color: var(--text-muted);">Hooks</p>
            </div>
            <div class="p-4 rounded-xl border" style="background-color: var(--bg-muted); border-color: var(--border-color);">
              <Bot size={20} style="color: var(--accent);" class="mb-2" />
              <p class="text-2xl font-bold" style="color: var(--text-primary);">3</p>
              <p class="text-sm" style="color: var(--text-muted);">Agents</p>
            </div>
            <div class="p-4 rounded-xl border" style="background-color: var(--bg-muted); border-color: var(--border-color);">
              <Server size={20} style="color: var(--accent);" class="mb-2" />
              <p class="text-2xl font-bold" style="color: var(--text-primary);">2</p>
              <p class="text-sm" style="color: var(--text-muted);">MCP Servers</p>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <button class="card p-5 flex items-center gap-4 hover:border-[var(--accent)] transition-all duration-200 cursor-pointer">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent-muted);">
              <Copy size={20} style="color: var(--accent);" />
            </div>
            <div class="text-left">
              <p class="font-medium" style="color: var(--text-primary);">Copy Entity</p>
              <p class="text-sm" style="color: var(--text-muted);">Copy to project</p>
            </div>
          </button>
          <button class="card p-5 flex items-center gap-4 hover:border-[var(--accent)] transition-all duration-200 cursor-pointer">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent-muted);">
              <GitCompare size={20} style="color: var(--accent);" />
            </div>
            <div class="text-left">
              <p class="font-medium" style="color: var(--text-primary);">Diff Skills</p>
              <p class="text-sm" style="color: var(--text-muted);">Compare projects</p>
            </div>
          </button>
          <button class="card p-5 flex items-center gap-4 hover:border-[var(--accent)] transition-all duration-200 cursor-pointer">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color: var(--accent-muted);">
              <Trash2 size={20} style="color: var(--destructive);" />
            </div>
            <div class="text-left">
              <p class="font-medium" style="color: var(--text-primary);">Clean Teams</p>
              <p class="text-sm" style="color: var(--text-muted);">Remove stale teams</p>
            </div>
          </button>
        </div>
      </div>

    {:else}
      <!-- Entities View -->
      <div class="space-y-6">
        <!-- Filter Bar -->
        <div class="flex items-center gap-2 overflow-x-auto pb-2">
          {#each entityFilters as filter}
            <button
              onclick={() => activeEntityFilter = filter}
              class="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-all duration-200"
              style="background-color: {activeEntityFilter === filter ? 'var(--accent)' : 'var(--bg-muted)'}; color: {activeEntityFilter === filter ? 'white' : 'var(--text-secondary)'};"
            >
              {filter === 'all' ? 'All' : filter.charAt(0).toUpperCase() + filter.slice(1) + 's'}
            </button>
          {/each}
        </div>

        <!-- Entities List -->
        <div class="space-y-2">
          {#each filteredEntities as entity, i}
            {@const Icon = entityIcons[entity.type] || Package}
            <div class="entity-list-item animate-in" style="animation-delay: {i * 0.03}s;">
              <div class="entity-list-item-icon">
                <Icon size={18} />
              </div>
              <div class="entity-list-item-content">
                <p class="entity-list-item-title">{entity.name}</p>
                <p class="entity-list-item-meta">{entity.count} instances across {entity.location === 'global' ? 'global' : 'projects'}</p>
              </div>
              <span class="badge {entity.location === 'global' ? 'badge-brand' : 'badge-success'}">
                {entity.location}
              </span>
              <div class="flex items-center gap-1">
                <button class="btn-icon w-8 h-8" aria-label="Copy">
                  <Copy size={14} />
                </button>
                <button class="btn-icon w-8 h-8" aria-label="View details">
                  <ChevronRight size={14} />
                </button>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </main>

  <!-- Footer -->
  <footer class="border-t mt-12 py-6" style="border-color: var(--border-color);">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row items-center justify-between gap-4">
        <p class="text-sm" style="color: var(--text-muted);">
          Built with care by <span class="font-medium" style="color: var(--accent);">@thieung</span>
        </p>
        <div class="flex items-center gap-4">
          <span class="badge badge-success flex items-center gap-1">
            <CheckCircle2 size={12} />
            All systems operational
          </span>
          <span class="text-xs font-mono" style="color: var(--text-muted);">v1.0.0</span>
        </div>
      </div>
    </div>
  </footer>
</div>
