<p align="center">
  <img src="https://img.shields.io/badge/🔍_ccinspect-Claude_Code_Inspector-8B5CF6?style=for-the-badge&labelColor=1e1e2e" alt="ccinspect" />
</p>

<p align="center">
  <strong>Scan, inventory, and manage all Claude Code configurations across your machine.</strong>
</p>

<p align="center">
  <a href="#install"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white" alt="Go 1.22+" /></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-green" alt="MIT License" /></a>
  <a href="https://github.com/thieung/ccinspect/releases"><img src="https://img.shields.io/github/v/release/thieung/ccinspect?color=purple" alt="Release" /></a>
  <a href="https://github.com/thieung/ccinspect/stargazers"><img src="https://img.shields.io/github/stars/thieung/ccinspect?style=social" alt="Stars" /></a>
</p>

---

## Key Features

- **Scan** — Find every `.claude/` directory (global + per-project) across configurable paths
- **Inventory** — Parse skills, hooks, agents, commands, rules, MCP servers, and teams
- **List** — Filter and browse entities by type with formatted tables or JSON
- **Copy** — Copy skills, agents, hooks, and commands between projects or from global config
- **Diff** — Compare skills between two projects side-by-side
- **Clean** — Remove `.claude/` config from a specific project (with dry-run support)
- **Sort** — Output sorted by entity count (most configured projects first)
- **Progress** — Animated spinner while scanning large directory trees

## Install

Works on **macOS**, **Linux**, and **Windows**.

### One-liner (recommended)

**macOS / Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/thieung/ccinspect/main/install.sh | sh
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/thieung/ccinspect/main/install.ps1 | iex
```

### With Go

```bash
go install github.com/thieung/ccinspect/cmd/ccinspect@latest
```

### From source

```bash
git clone https://github.com/thieung/ccinspect.git
cd ccinspect
make install    # builds + copies to /usr/local/bin
```

### Uninstall

```bash
# If installed via install script or make install
sudo rm /usr/local/bin/ccinspect

# Windows
Remove-Item "$env:LOCALAPPDATA\ccinspect" -Recurse
```

## Quick Start

```bash
# Initialize config (optional — uses sensible defaults)
ccinspect init

# Scan all projects
ccinspect scan

# Scan a single project
ccinspect scan --project ~/projects/my-app

# Global config summary
ccinspect global

# Clean a project's .claude/ directory
ccinspect clean ~/projects/old-app --dry-run
```

## Commands

### `ccinspect scan`

Summary table of all Claude Code installations found. Projects are sorted by total entity count (most configured first).

```bash
ccinspect scan              # all projects
ccinspect scan --project .  # current project only
ccinspect scan --json       # JSON output
```

### `ccinspect list <entity>`

List all entities of a given type.

```bash
ccinspect list skills                  # all skills (global + projects)
ccinspect list skills --global-only    # global skills only
ccinspect list hooks                   # all hooks
ccinspect list agents                  # all agents
ccinspect list commands                # all commands
ccinspect list rules                   # all rules
ccinspect list mcp                     # all MCP servers
ccinspect list teams                   # all teams
ccinspect list skills --json           # JSON output
```

### `ccinspect global`

Show global `~/.claude/` configuration summary with entity counts.

```bash
ccinspect global
ccinspect global --json
```

### `ccinspect diff skills <A> <B>`

Compare skills between two projects. Use `global` as alias for `~/.claude/`.

```bash
ccinspect diff skills global ~/projects/my-app
ccinspect diff skills ~/projects/app-a ~/projects/app-b
ccinspect diff skills global ~/projects/my-app --json
```

### `ccinspect copy <entity-type> <name> [name2 ...]`

Copy skills, agents, hooks, or commands between projects or between global (`~/.claude`) and project configs.

**Defaults:** `--from=global`, `--to=.` (current directory) — so the simplest usage is just:

```bash
# Copy a skill from global → current project (shortest form)
ccinspect copy skills my-skill

# Copy all skills from global → current project
ccinspect copy skills all

# Copy multiple agents from global → current project
ccinspect copy agents planner debugger
```

More examples with explicit flags:

```bash
# Copy from global to a specific project
ccinspect copy skills my-skill --to ~/projects/my-app

# Copy between two projects
ccinspect copy skills all --from ~/projects/app-a --to ~/projects/app-b

# Copy hook events
ccinspect copy hooks PreToolUse PostToolUse --to ~/projects/my-app

# Copy all commands between projects
ccinspect copy commands all --from ~/projects/app-a --to ~/projects/app-b

# Preview without making changes
ccinspect copy skills all --dry-run

# Force overwrite existing
ccinspect copy skills my-skill --force

# List available entities in a source
ccinspect copy skills all --list
ccinspect copy hooks all --list

# JSON output
ccinspect copy skills my-skill --json
```

| Flag | Default | Description |
|------|---------|-------------|
| `--from` | `global` | Source: `global` or project path |
| `--to` | `.` | Destination: `global` or project path |
| `--force` | | Overwrite existing entities at destination |
| `--dry-run` | | Show what would be copied without making changes |
| `--list` | | List available entities in the source |
| `--json` | | Output as JSON |

### `ccinspect clean <project-path>`

Remove `.claude/` directory (and `.mcp.json`) from a project.

```bash
ccinspect clean ~/projects/old-app --dry-run   # preview what will be removed
ccinspect clean ~/projects/old-app              # actually remove
ccinspect clean .                               # clean current directory

# Clean stale teams (no config.json — runtime/temporary teams)
ccinspect clean teams --dry-run                 # preview
ccinspect clean teams                           # remove stale teams only
ccinspect clean teams --all                     # remove ALL teams
```

### `ccinspect config show`

Display current configuration.

### `ccinspect init`

Create default config at `~/.ccinspect/config.json`.

## Configuration

Config file: `~/.ccinspect/config.json`

```json
{
  "scan_paths": ["~/projects", "~/work"],
  "exclude_paths": ["node_modules", ".git", "vendor"],
  "max_depth": 5,
  "default_output": "table"
}
```

## Entities Scanned

| Entity | Location | Source |
|--------|----------|--------|
| Skills | `skills/` subdirs with `SKILL.md` | Global + Project |
| Hooks | `hooks` key in `settings.json` | Global + Project |
| Agents | `.md` files in `agents/` | Global |
| Commands | `.md` files in `commands/` | Global + Project |
| Rules | `.md` files in `rules/` | Global |
| MCP Servers | `.mcp.json`, `settings.json`, `settings.local.json` | Project |
| Teams | `teams/` subdirs with `config.json` | Global |

## Roadmap

- [x] Scan & inventory all Claude Code installations
- [x] List entities by type (skills, hooks, agents, MCP, etc.)
- [x] Diff skills between projects
- [x] Clean command with dry-run
- [x] Copy entities between projects and global config
- [x] Sorted output by entity count
- [x] Progress spinner animation
- [x] Nested project scanning
- [ ] Export inventory to HTML report
- [ ] Skill sync & versioning (push, divergence detection)
- [ ] Session log viewer (list, show, live tail)
- [ ] Web UI dashboard with embedded Svelte frontend
- [ ] Hallucination detection & error pattern analysis
- [ ] Watch mode — auto-rescan on config changes
- [ ] LLM-powered usage insights & suggestions

## Contributing

Contributions are welcome! Please open an issue first to discuss what you'd like to change.

```bash
# Development
git clone https://github.com/thieung/ccinspect.git
cd ccinspect
make build    # binary at ./bin/ccinspect
make test
```

1. Fork the repo
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

## Star History

<p align="center">
  <a href="https://star-history.com/#thieung/ccinspect&Date">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=thieung/ccinspect&type=Date&theme=dark" />
      <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=thieung/ccinspect&type=Date" />
      <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=thieung/ccinspect&type=Date" />
    </picture>
  </a>
</p>

## ❤️ Support This Project

This project is used by developers to scan, inventory, and manage Claude Code configurations across machines.

Maintaining it takes time — reviewing issues, fixing bugs, and building new features.

If it helps your work, consider supporting:

👉 GitHub Sponsors: https://github.com/sponsors/thieung

Your support helps me:
- Ship features faster
- Maintain long-term stability
- Keep this project open and free

## License

[MIT](LICENSE) — made with care by [@thieung](https://github.com/thieung)
