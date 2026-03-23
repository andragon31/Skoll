// @ts-nocheck
/* eslint-disable */
declare var Bun: any;
declare var process: any;

import type { Plugin } from "@opencode-ai/plugin"

// ─── Configuration ───────────────────────────────────────────────────────────

const SKOLL_BIN = process.env.SKOLL_BIN ?? "skoll"

const SKOLL_TOOLS = new Set([
  "rsaw_scan",
  "rsaw_read_item",
  "rsaw_create_item",
  "rsaw_update_item",
  "rsaw_get_template",
  "rule_list",
])

// ─── Orchestration Instructions ──────────────────────────────────────────────

const ORCHESTRATION_INSTRUCTIONS = `## Skoll Protocol
You have access to Skoll, an AI Orchestration Layer for the RSAW (Rules, Skills, Agents, Workflows) framework.

### ORCHESTRATION TOOLS:
- **rsaw_scan()**: MANDATORY at session start to see all project components.
- **rsaw_read_item(path)**: Use to analyze existing workflows, agents or skills.
- **rsaw_create_item(type, name, content)**: Use to generate new RSAW components dynamically.
- **rsaw_update_item(path, content)**: Use to link agents to workflows or update rules.
- **rsaw_get_template(type)**: Get the base structure for a component.

### RULES:
1. Always scan the project with \`rsaw_scan\` before designing new components.
2. Follow the templates obtained via \`rsaw_get_template\`.
3. Use \`rsaw_update_item\` to maintain the linkages between agents and workflows.
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
