package main

import (
	"fmt"

	"github.com/andragon31/skoll/internal/analyzer"
	"github.com/andragon31/skoll/internal/generator"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze project stack & architecture to generate RSAW components.",
	RunE: func(cmd *cobra.Command, args []string) error {
		meta := analyzer.Analyze(".")
		fmt.Printf("Skoll Project analysis:\nLanguage: %s\nBuild Tools: %v\nModules found: %v\n", meta.Language, meta.BuildTools, meta.Modules)

		fmt.Println("\nCreating default RSAW structure...")
		return generator.InitProjectWithMeta(".", meta)
	},
}
