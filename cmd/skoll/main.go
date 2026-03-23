package main

import (
	"embed"
	"fmt"
	"runtime"

	"github.com/andragon31/skoll/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:embed plugins/opencode/*
var openCodeFS embed.FS

var (
	version = "0.1.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "skoll",
		Short: "Skoll - RSAW Orchestration Layer for AI Agents",
		Long: `Skoll is a structured orchestration layer for AI development teams.
It implements the RSAW (Rules, Skills, Agents, Workflows) framework to bring 
consistency, clear roles, and reproducible processes to your AI-assisted workflow.`,
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("data-dir", "d", "", "Data directory (default: ~/.skoll)")
	rootCmd.PersistentFlags().String("log-level", "info", "Log level: debug|info|warn|error")

	viper.BindPFlag("data_dir", rootCmd.PersistentFlags().Lookup("data-dir"))
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Skoll v%s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
		},
	}
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(tuiCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

func initConfig() {
	viper.SetDefault("data_dir", utils.GetDefaultDataDir())
	viper.SetDefault("log_level", "info")
}
