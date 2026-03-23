// @ts-nocheck
/* eslint-disable */
declare var Bun: any;
declare var process: any;

import type { Plugin } from "@opencode-ai/plugin"

// ─── Configuration ───────────────────────────────────────────────────────────

const SKOLL_PORT = parseInt(Bun.env.SKOLL_PORT ?? process.env.SKOLL_PORT ?? "7439")
const SKOLL_URL = `http://127.0.0.1:${SKOLL_PORT}`
const SKOLL_BIN = process.env.SKOLL_BIN ?? "skoll"

const SKOLL_TOOLS = new Set([
  "rule_list",
  "rule_check",
  "skill_load",
  "skill_search",
  "agent_activate",
  "agent_list",
  "agent_context",
  "agent_handoff",
  "workflow_start",
  "workflow_step",
  "workflow_status",
  "workflow_complete",
])

// ─── Orchestration Instructions ──────────────────────────────────────────────

const ORCHESTRATION_INSTRUCTIONS = `## Skoll Protocol
You have access to Skoll, a RSAW (Rules, Skills, Agents, Workflows) Orchestration Layer.

### ORCHESTRATION TOOLS:
#### Rules (MANDATORY at start)
Call: rule_list() to understand global project constraints.

#### Agents (MANDATORY when switching context)
Call: agent_activate(name="agent_name") to load identity, scope and skills for a specific role (e.g. "backend", "frontend").

#### Skills (MANDATORY before performing technical tasks)
Call: skill_load(name="skill_name") to load specific project "how-to" knowledge (e.g. "git-workflow").

#### Workflows (MANDATORY for complex multi-step processes)
Call: workflow_start(name="workflow_name") to initiate a guided process.
`

export const Skoll: Plugin = async (ctx) => {
  return {
    "experimental.chat.system.transform": async (_input, output) => {
      if (output.system.length > 0) {
        output.system[output.system.length - 1] += "\n\n" + ORCHESTRATION_INSTRUCTIONS
      } else {
        output.system.push(ORCHESTRATION_INSTRUCTIONS)
      }
    },
  }
}
