package main

import (
	"github.com/andragon31/skoll/internal/mcp"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start Skoll MCP server (stdio)",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.Default()
		srv := mcp.NewServer(logger)
		return srv.RunStdio()
	},
}
