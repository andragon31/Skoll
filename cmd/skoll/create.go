package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	createType    string
	createName    string
	createContent string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new RSAW component (Agent, Rule, Skill, Workflow)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if createType == "" || createName == "" || createContent == "" {
			return fmt.Errorf("flags --type, --name and --content are required")
		}

		// Ensure directory exists
		dir := filepath.Join(".skoll", strings.ToLower(createType)+"s")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		destPath := filepath.Join(dir, strings.ToLower(createName)+".md")
		fmt.Printf("Skoll creating %s...\n", destPath)
		return os.WriteFile(destPath, []byte(createContent), 0644)
	},
}

func init() {
	createCmd.Flags().StringVarP(&createType, "type", "t", "", "Component type (Agent, Rule, Skill, Workflow)")
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Component name")
	createCmd.Flags().StringVarP(&createContent, "content", "c", "", "Component content (Markdown)")
}
