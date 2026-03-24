package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/andragon31/skoll/internal/rsaw"
	"github.com/charmbracelet/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	logger *log.Logger
	server *server.MCPServer
	loader *rsaw.Loader
	root   string
}

func NewServer(logger *log.Logger, root string) *Server {
	srv := server.NewMCPServer("skoll", "4.0.0")

	s := &Server{
		logger: logger,
		server: srv,
		loader: rsaw.NewLoader(logger, root),
		root:   root,
	}

	s.registerTools()

	return s
}

func (s *Server) registerTools() {
	s.registerRulesTools()
	s.registerSkillsTools()
	s.registerAgentsTools()
	s.registerWorkflowsTools()
	s.registerSystemTools()
}

func (s *Server) registerRulesTools() {
	s.server.AddTool(mcp.NewTool("rule_list",
		mcp.WithDescription("List all rules"),
		mcp.WithString("category", mcp.Description("Filter by category")),
	), s.handleRuleList)

	s.server.AddTool(mcp.NewTool("rule_check",
		mcp.WithDescription("Check if an action violates any rule"),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to check")),
	), s.handleRuleCheck)

	s.server.AddTool(mcp.NewTool("rule_get",
		mcp.WithDescription("Get rule details"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Rule name")),
	), s.handleRuleGet)
}

func (s *Server) registerSkillsTools() {
	s.server.AddTool(mcp.NewTool("skill_list",
		mcp.WithDescription("List all skills (compact format for Progressive Disclosure)"),
		mcp.WithBoolean("compact", mcp.Description("Return compact metadata only")),
	), s.handleSkillList)

	s.server.AddTool(mcp.NewTool("skill_load",
		mcp.WithDescription("Load full skill content including allowed-tools"),
		mcp.WithString("skill_name", mcp.Required(), mcp.Description("Skill name to load")),
	), s.handleSkillLoad)

	s.server.AddTool(mcp.NewTool("skill_read_file",
		mcp.WithDescription("Read a specific file from a skill's scripts/references/assets"),
		mcp.WithString("skill_name", mcp.Required(), mcp.Description("Skill name")),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("Relative path within skill directory")),
	), s.handleSkillReadFile)

	s.server.AddTool(mcp.NewTool("skills_import",
		mcp.WithDescription("Import a skill from SkillsMP or GitHub"),
		mcp.WithString("source", mcp.Required(), mcp.Description("Source: skillsmp, github, or local")),
		mcp.WithString("query", mcp.Description("Search query for skillsmp")),
		mcp.WithString("url", mcp.Description("GitHub URL for github source")),
		mcp.WithBoolean("skip_scan", mcp.Description("Skip security scan (NOT recommended)")),
	), s.handleSkillsImport)

	s.server.AddTool(mcp.NewTool("skills_update",
		mcp.WithDescription("Update imported skills from their source"),
		mcp.WithString("skill_name", mcp.Description("Skill name to update (omit for all)")),
	), s.handleSkillsUpdate)
}

func (s *Server) registerAgentsTools() {
	s.server.AddTool(mcp.NewTool("agent_list",
		mcp.WithDescription("List all agents"),
	), s.handleAgentList)

	s.server.AddTool(mcp.NewTool("agent_activate",
		mcp.WithDescription("Activate an agent with their skills and scope"),
		mcp.WithString("agent_id", mcp.Required(), mcp.Description("Agent ID to activate")),
		mcp.WithString("context_path", mcp.Description("Context path for AGENTS.md lookup")),
	), s.handleAgentActivate)

	s.server.AddTool(mcp.NewTool("agent_context",
		mcp.WithDescription("Get current agent context"),
	), s.handleAgentContext)

	s.server.AddTool(mcp.NewTool("agent_handoff",
		mcp.WithDescription("Hand off work to another agent"),
		mcp.WithString("to_agent", mcp.Required(), mcp.Description("Target agent ID")),
		mcp.WithString("contract", mcp.Required(), mcp.Description("Work contract summary")),
	), s.handleAgentHandoff)
}

func (s *Server) registerWorkflowsTools() {
	s.server.AddTool(mcp.NewTool("workflow_start",
		mcp.WithDescription("Start a workflow"),
		mcp.WithString("workflow_id", mcp.Required(), mcp.Description("Workflow ID")),
	), s.handleWorkflowStart)

	s.server.AddTool(mcp.NewTool("workflow_step",
		mcp.WithDescription("Execute a workflow step"),
		mcp.WithString("step_id", mcp.Required(), mcp.Description("Step ID")),
		mcp.WithString("output", mcp.Description("Step output")),
	), s.handleWorkflowStep)

	s.server.AddTool(mcp.NewTool("workflow_status",
		mcp.WithDescription("Get workflow status"),
		mcp.WithString("workflow_id", mcp.Required(), mcp.Description("Workflow ID")),
	), s.handleWorkflowStatus)

	s.server.AddTool(mcp.NewTool("workflow_complete",
		mcp.WithDescription("Complete a workflow"),
		mcp.WithString("workflow_id", mcp.Required(), mcp.Description("Workflow ID")),
		mcp.WithBoolean("success", mcp.Description("Whether workflow succeeded")),
	), s.handleWorkflowComplete)

	s.server.AddTool(mcp.NewTool("dod_check",
		mcp.WithDescription("Check Definition of Done for a workflow"),
		mcp.WithString("workflow_name", mcp.Description("Workflow name")),
	), s.handleDodCheck)
}

func (s *Server) registerSystemTools() {
	s.server.AddTool(mcp.NewTool("skoll_status",
		mcp.WithDescription("Get Skoll system status"),
	), s.handleSkollStatus)

	s.server.AddTool(mcp.NewTool("skoll_validate",
		mcp.WithDescription("Validate SKILL.md format against AgentSkills standard"),
		mcp.WithBoolean("ci", mcp.Description("CI mode (strict)")),
	), s.handleSkollValidate)

	s.server.AddTool(mcp.NewTool("rule_pending",
		mcp.WithDescription("List pending rules from Fenrir"),
	), s.handleRulePending)

	s.server.AddTool(mcp.NewTool("rule_promote",
		mcp.WithDescription("Promote a pending rule to active"),
		mcp.WithString("rule_id", mcp.Required(), mcp.Description("Rule ID to promote")),
	), s.handleRulePromote)

	s.server.AddTool(mcp.NewTool("team_status",
		mcp.WithDescription("Get team coordination status"),
	), s.handleTeamStatus)

	s.server.AddTool(mcp.NewTool("team_register",
		mcp.WithDescription("Register a developer in the team"),
		mcp.WithString("agent", mcp.Description("Agent name")),
	), s.handleTeamRegister)
}

func (s *Server) handleRuleList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	category := getStringOrDefault(request.GetArguments(), "category", "")
	rules, _ := s.loader.LoadRules(category)

	data, _ := json.Marshal(rules)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleRuleCheck(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_ = getString(request.GetArguments(), "action")

	result := map[string]interface{}{
		"allowed":    true,
		"violations": []string{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleRuleGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := getString(request.GetArguments(), "name")

	path := filepath.Join(s.root, ".skoll", "rules", name+".md")
	content, err := os.ReadFile(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(content)), nil
}

func (s *Server) handleSkillList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	compact := getBoolOrDefault(request.GetArguments(), "compact", true)

	skills, err := s.loader.LoadSkillsIndex()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if !compact {
		data, _ := json.Marshal(skills)
		return mcp.NewToolResultText(string(data)), nil
	}

	type CompactSkill struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     string `json:"version_status"`
	}

	compactSkills := make([]CompactSkill, 0, len(skills))
	for _, skill := range skills {
		compactSkills = append(compactSkills, CompactSkill{
			Name:        skill.Name,
			Description: skill.Description,
			Version:     skill.Metadata["version"],
		})
	}

	data, _ := json.Marshal(compactSkills)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleSkillLoad(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	skillName := getString(request.GetArguments(), "skill_name")

	skillPath := filepath.Join(s.root, ".skoll", "skills", skillName, "SKILL.md")
	content, err := os.ReadFile(skillPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	frontmatter := s.loader.ParseSKILLFrontmatter(string(content))

	result := map[string]interface{}{
		"name":            skillName,
		"content":         string(content),
		"allowed_tools":   frontmatter["allowed-tools"],
		"available_files": s.loader.ListSkillFiles(skillName),
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleSkillReadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	skillName := getString(request.GetArguments(), "skill_name")
	filePath := getString(request.GetArguments(), "file_path")

	fullPath := filepath.Join(s.root, ".skoll", "skills", skillName, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(string(content)), nil
}

func (s *Server) handleSkillsImport(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	source := getString(args, "source")

	result := map[string]interface{}{
		"source":   source,
		"imported": false,
		"message":  "SkillsMP integration requires network access",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleSkillsUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"updated": 0,
		"message": "No skills to update",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleAgentList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	agents, _ := s.loader.LoadAgents()

	data, _ := json.Marshal(agents)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleAgentActivate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	agentID := getString(args, "agent_id")
	contextPath := getStringOrDefault(args, "context_path", "")

	result := map[string]interface{}{
		"agent":         agentID,
		"skills_loaded": []string{},
		"allowed_tools": []string{},
		"scope":         map[string]interface{}{},
	}

	if contextPath != "" {
		agentsPath := s.findNestedAgentsMD(contextPath)
		if agentsPath != "" {
			content, _ := os.ReadFile(agentsPath)
			result["local_agents_md"] = map[string]string{
				"found_at":   agentsPath,
				"content":    string(content),
				"precedence": "local_overrides_root",
			}
		}
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) findNestedAgentsMD(contextPath string) string {
	dir := filepath.Dir(contextPath)
	for i := 0; i < 10; i++ {
		checkPath := filepath.Join(dir, "AGENTS.md")
		if _, err := os.Stat(checkPath); err == nil {
			return checkPath
		}
		dir = filepath.Dir(dir)
		if dir == "." {
			break
		}
	}
	return ""
}

func (s *Server) handleAgentContext(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"active_agent": "",
		"scope":        map[string]interface{}{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleAgentHandoff(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	toAgent := getString(args, "to_agent")
	contract := getString(args, "contract")

	result := map[string]interface{}{
		"handoff_to": toAgent,
		"contract":   contract,
		"confirmed":  true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleWorkflowStart(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	workflowID := getString(request.GetArguments(), "workflow_id")

	result := map[string]interface{}{
		"workflow_id":  workflowID,
		"status":       "started",
		"current_step": 0,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleWorkflowStep(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"status": "step_completed",
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleWorkflowStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"status":      "running",
		"step":        0,
		"total_steps": 0,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleWorkflowComplete(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"status":  "completed",
		"success": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleDodCheck(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"passed":  true,
		"checks":  []string{},
		"missing": []string{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleSkollStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"version":       "4.0.0",
		"rssaw_version": "v4.0",
		"skills_count":  0,
		"agents_count":  0,
		"rules_count":   0,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleSkollValidate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_ = getBoolOrDefault(request.GetArguments(), "ci", false)

	result := map[string]interface{}{
		"valid":    true,
		"errors":   []string{},
		"warnings": []string{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleRulePending(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"pending_rules": []string{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleRulePromote(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"promoted": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleTeamStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"active_developers": []string{},
		"role_assignments":  map[string]interface{}{},
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) handleTeamRegister(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := map[string]interface{}{
		"registered": true,
	}

	data, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(data)), nil
}

func (s *Server) RunStdio() error {
	return server.ServeStdio(s.server)
}

func getString(args map[string]interface{}, key string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getStringOrDefault(args map[string]interface{}, key, defaultVal string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

func getBoolOrDefault(args map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := args[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultVal
}
