package generator

const SkollProtocolMarkdown = `## Skoll Protocol
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
