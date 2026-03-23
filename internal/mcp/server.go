package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andragon31/skoll/internal/generator"
	"github.com/andragon31/skoll/internal/rsaw"
	"github.com/charmbracelet/log"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	logger *log.Logger
	server *server.MCPServer
	loader *rsaw.Loader
}

func NewServer(logger *log.Logger) *Server {
	srv := server.NewMCPServer("skoll", "0.1.0")

	s := &Server{
		logger: logger,
		server: srv,
		loader: rsaw.NewLoader(logger),
	}

	s.registerTools()

	return s
}

func (s *Server) registerTools() {
	// RSAW Scan & List
	s.server.AddTool(mcp.NewTool("rsaw_scan",
		mcp.WithDescription("Scans and lists all Rules, Skills, Agents and Workflows in the current project"),
	), s.handleScan)

	// RSAW Read
	s.server.AddTool(mcp.NewTool("rsaw_read_item",
		mcp.WithDescription("Reads the content of an RSAW item by Its path"),
		mcp.WithString("path", mcp.Required(), mcp.Description("Full path to the item file")),
	), s.handleReadItem)

	// RSAW Create
	s.server.AddTool(mcp.NewTool("rsaw_create_item",
		mcp.WithDescription("Creates a new RSAW item"),
		mcp.WithString("type", mcp.Required(), mcp.Description("Component type (rule, skill, agent, workflow)")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name of the component")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Markdown content following the template")),
	), s.handleCreateItem)

	// RSAW Update
	s.server.AddTool(mcp.NewTool("rsaw_update_item",
		mcp.WithDescription("Updates an existing RSAW item's content"),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to the file to update")),
		mcp.WithString("content", mcp.Required(), mcp.Description("New content for the file")),
	), s.handleUpdateItem)

	// RSAW Get Template
	s.server.AddTool(mcp.NewTool("rsaw_get_template",
		mcp.WithDescription("Gets the structure template for a Rule, Skill, Agent or Workflow"),
		mcp.WithString("type", mcp.Required(), mcp.Description("Component type (rule, skill, agent, workflow)")),
	), s.handleGetTemplate)
}

func (s *Server) handleScan(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	items, err := s.loader.ScanProject()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Scanning project: %v", err)), nil
	}
	res, _ := json.Marshal(items)
	return mcp.NewToolResultText(string(res)), nil
}

func (s *Server) handleReadItem(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}
	path := args["path"].(string)
	content, err := os.ReadFile(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error reading file: %v", err)), nil
	}
	return mcp.NewToolResultText(string(content)), nil
}

func (s *Server) handleCreateItem(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}
	typeName := strings.ToLower(args["type"].(string))
	name := args["name"].(string)
	content := args["content"].(string)

	dir := filepath.Join(".skoll", typeName+"s")
	os.MkdirAll(dir, 0755)

	path := filepath.Join(dir, strings.ToLower(name)+".md")
	os.WriteFile(path, []byte(content), 0644)
	return mcp.NewToolResultText(fmt.Sprintf("Item created at: %s", path)), nil
}

func (s *Server) handleUpdateItem(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}
	path := args["path"].(string)
	content := args["content"].(string)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error updating file: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Item at %s updated successfully", path)), nil
}

func (s *Server) handleGetTemplate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments"), nil
	}
	typeName := strings.ToLower(args["type"].(string))
	content, err := generator.GetTemplateContent(typeName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Template not found: %v", err)), nil
	}
	return mcp.NewToolResultText(content), nil
}

func (s *Server) RunStdio() error {
	return server.ServeStdio(s.server)
}
