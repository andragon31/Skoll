package main

import (
	"os"

	"github.com/andragon31/skoll/internal/mcp"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start Skoll MCP server (stdio)",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.Default()
		root, _ := os.Getwd()
		srv := mcp.NewServer(logger, root)
		return srv.RunStdio()
	},
}
