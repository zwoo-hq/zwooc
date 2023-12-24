package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	pendingSpinner = spinner.Spinner{
		Frames: []string{
			"▰▱▱▱▱",
			"▱▰▱▱▱",
			"▱▱▰▱▱",
			"▱▱▱▰▱",
			"▱▱▱▱▰",
			"▱▱▱▰▱",
			"▱▱▰▱▱",
			"▱▰▱▱▱",
			"▰▱▱▱▱",
		},
		FPS: time.Second / 6,
	}

	runningSpinner = spinner.Spinner{
		Frames: []string{
			"▱▱▱▱▱",
			"▰▱▱▱▱",
			"▰▰▱▱▱",
			"▰▰▰▱▱",
			"▰▰▰▰▰",
			"▱▰▰▰▰",
			"▱▱▰▰▰",
			"▱▱▱▰▰",
			"▱▱▱▱▰",
			"▱▱▱▱▱",
		},
		FPS: time.Second / 8,
	}

	pendingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	runningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Padding(0, 2).SetString("✓")
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("124")).Padding(0, 2).SetString("✗")
	canceledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Padding(0, 2).SetString("-")
)

func PrintSuccess(name string, dur time.Duration) {
	fmt.Printf("%s successfully ran %s in %s\n", successStyle.Render(), name, dur)
}

func PrintError(name string) {
	fmt.Printf("%s failed running %s\n", errorStyle.Render(), name)
}

func PrintCancel(name string) {
	fmt.Printf("%s %s was canceled\n", canceledStyle.Render(), name)
}
