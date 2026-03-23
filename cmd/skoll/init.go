package main

import (
	"github.com/andragon31/skoll/internal/generator"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Skoll in current project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return generator.InitProject(".")
	},
}

func init() {
	// Register command in main (I'll update main.go after this)
}
