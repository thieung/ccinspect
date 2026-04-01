# ccinspect

CLI tool to scan and inventory all Claude Code configurations across your machine.

Finds every `.claude/` directory (global + per-project), parses config entities (skills, hooks, agents, commands, rules, MCP servers, teams), and outputs a summary table.

## Install

```bash
# From source
go install github.com/thieunv/ccinspect/cmd/cli@latest

# Or build locally
git clone https://github.com/thieunv/ccinspect.git
cd ccinspect
make build
# Binary at ./bin/ccinspect
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
```

## Commands

### `ccinspect scan`

Summary table of all Claude Code installations found.

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
| MCP Servers | `.mcp.json` at project root | Project |
| Teams | `teams/` subdirs with `config.json` | Global |

## Requirements

- Go 1.22+
