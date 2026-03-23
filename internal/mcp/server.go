package mcp

import (
	"context"
	"fmt"

	"github.com/andragon31/skoll/internal/rsaw"
	"github.com/charmbracelet/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	loader *rsaw.Loader
	logger *log.Logger
	server *server.MCPServer
}

func NewServer(l *rsaw.Loader, logger *log.Logger) *Server {
	srv := server.NewMCPServer("skoll", "0.1.0")

	s := &Server{
		loader: l,
		logger: logger,
		server: srv,
	}

	s.registerAllTools()

	return s
}

func (s *Server) registerAllTools() {
	s.registerRuleTools()
	s.registerSkillTools()
	s.registerAgentTools()
	s.registerWorkflowTools()
}

func (s *Server) registerRuleTools() {
	ruleList := mcp.NewTool("rule_list",
		mcp.WithDescription("List all active rules in the project"),
	)
	s.server.AddTool(ruleList, s.handleRuleList)
}

func (s *Server) registerSkillTools() {
	skillLoad := mcp.NewTool("skill_load",
		mcp.WithDescription("Load a specific skill by name"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Skill name (e.g. git-workflow)")),
	)
	s.server.AddTool(skillLoad, s.handleSkillLoad)
}

func (s *Server) registerAgentTools() {
	agentActivate := mcp.NewTool("agent_activate",
		mcp.WithDescription("Activate an agent persona with its scope and skills"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Agent name")),
	)
	s.server.AddTool(agentActivate, s.handleAgentActivate)
}

func (s *Server) registerWorkflowTools() {
	workflowStart := mcp.NewTool("workflow_start",
		mcp.WithDescription("Start a new structured workflow"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Workflow name")),
	)
	s.server.AddTool(workflowStart, s.handleWorkflowStart)
}

// Handlers (stubs for now)
func (s *Server) handleRuleList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	items, err := s.loader.ScanProject()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var rules []string
	for _, item := range items {
		if item.Type == rsaw.TypeRule {
			rules = append(rules, item.Name)
		}
	}
	return mcp.NewToolResultText(fmt.Sprintf("Rules found: %v", rules)), nil
}

func (s *Server) handleSkillLoad(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("Skill loaded (stub)"), nil
}

func (s *Server) handleAgentActivate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("Agent activated (stub)"), nil
}

func (s *Server) handleWorkflowStart(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText("Workflow started (stub)"), nil
}

func (s *Server) RunStdio() error {
	return server.ServeStdio(s.server)
}
