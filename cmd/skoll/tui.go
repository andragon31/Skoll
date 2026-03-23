package main

import (
	"github.com/andragon31/skoll/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Open Skoll Interactive Dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
		_, err := p.Run()
		return err
	},
}
