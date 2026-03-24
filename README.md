# Skoll

**RSAW Orchestration Layer for AI Development Teams**

<p align="center">
<em>Rules, Skills, Agents, Workflows - Structured AI Orchestration</em>
</p>

Skoll provides a structured framework for organizing AI agent behaviors through Rules, Skills, Agents, and Workflows.

```
OpenCode / Claude Code / Cursor / Windsurf / ...
    ↓ MCP stdio
Skoll (single Go binary)
    ↓
.skoll/ directory in your project
```

## Features

- **Progressive Disclosure** - Compact skill listings, full content on demand
- **SKILL.md Format** - YAML frontmatter with allowed-tools specification
- **AGENTS.md Support** - Nested agent configs in monorepos
- **Team Coordination** - Register and track developer agents
- **Workflows** - Structured multi-step processes with DoD checks

## Quick Start

### Install (One-liner)

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/andragon31/Skoll/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/andragon31/Skoll/main/install.ps1 | iex
```

### Initialize in Your Project

```bash
skoll init
```

This creates the `.skoll/` directory structure:
```
.skoll/
├── rules/          # Team rules and guidelines
├── skills/        # Reusable skill definitions
├── agents/        # Agent configurations
└── workflows/     # Workflow definitions
```

### Setup Your Agent

| Agent | Command |
|-------|---------|
| OpenCode | `skoll setup opencode` |
| Claude Code | `skoll setup claude-code` |
| Cursor | `skoll setup cursor` |
| Windsurf | `skoll setup windsurf` |
| Antigravity | `skoll setup antigravity` |
| Gemini CLI | `skoll setup gemini-cli` |

## MCP Tools (21 total)

### Rules
| Tool | Description |
|------|-------------|
| `rule_list` | List all rules (with optional category filter) |
| `rule_check` | Check if an action violates any rule |
| `rule_get` | Get rule details by name |

### Skills
| Tool | Description |
|------|-------------|
| `skill_list` | List all skills (compact metadata by default) |
| `skill_load` | Load full skill content including allowed-tools |
| `skill_read_file` | Read file from skill's scripts/references/assets |
| `skills_import` | Import skill from SkillsMP or GitHub |
| `skills_update` | Update imported skills from source |

### Agents
| Tool | Description |
|------|-------------|
| `agent_list` | List all agents |
| `agent_activate` | Activate agent with skills and scope |
| `agent_context` | Get current agent context |
| `agent_handoff` | Hand off work to another agent |

### Workflows
| Tool | Description |
|------|-------------|
| `workflow_start` | Start a workflow |
| `workflow_step` | Execute a workflow step |
| `workflow_status` | Get workflow status |
| `workflow_complete` | Complete a workflow |
| `dod_check` | Check Definition of Done for workflow |

### System
| Tool | Description |
|------|-------------|
| `skoll_status` | Get Skoll system status |
| `skoll_validate` | Validate SKILL.md format |
| `rule_pending` | List pending rules |
| `rule_promote` | Promote pending rule to active |
| `team_status` | Get team coordination status |
| `team_register` | Register developer in team |

## CLI Reference

```bash
skoll setup [agent]   # Setup for an AI agent
skoll init           # Initialize in project
skoll mcp            # Start MCP server
skoll tui            # Open Dashboard
skoll version        # Show version
```

## SKILL.md Format

```yaml
---
name: my-skill
description: |
  What this skill does.
license: MIT
metadata:
  author: team
  version: "1.0"
allowed-tools:
  - mem_save
  - mem_find
---

## Cuándo aplicar

When to use this skill.

## Proceso

Step-by-step process.

## Checklist

- [ ] Step 1
- [ ] Step 2

## Anti-patrones

What NOT to do.
```

## Architecture

```
┌─────────────────────────────────────────────┐
│                 OpenCode                     │
│              Claude Code                     │
│                Cursor                        │
└─────────────────┬───────────────────────────┘
                  │ MCP stdio
                  ▼
┌─────────────────────────────────────────────┐
│                   Skoll                      │
├─────────────────────────────────────────────┤
│  Rules    │  Skills   │  Agents │ Workflows│
└─────────────────┬───────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────┐
│              .skoll/ directory              │
│  (rules/, skills/, agents/, workflows/)    │
└─────────────────────────────────────────────┘
```

## Documentation

- Create rules in `.skoll/rules/`
- Create skills in `.skoll/skills/` with SKILL.md
- Define agents in `.skoll/agents/`
- Define workflows in `.skoll/workflows/`

## License

MIT
