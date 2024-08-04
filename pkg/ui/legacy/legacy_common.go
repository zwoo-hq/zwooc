package legacyui

import (
	"fmt"
	"os"
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

	shutdownTabSpinner = spinner.Spinner{
		Frames: []string{
			"⣿⣿",
			"⣾⣷",
			"⣶⣶",
			"⣴⣦",
			"⣤⣤",
			"⣠⣄",
			"⣀⣀",
			"⣄⣠",
			"⣤⣤",
			"⣦⣴",
			"⣶⣶",
			"⣷⣾",
		},
		FPS: time.Second / 10,
	}

	pendingTabSpinner = spinner.Spinner{
		Frames: []string{
			"⡇ ",
			"⢸ ",
			"⢸⡇",
			" ⣿",
			" ⢸",
			" ⣹",
			" ⣉",
			"⣉ ",
			"⣏ ",
			"⡇ ",
		},
		FPS: time.Second / 10,
	}

	runningTabSpinner = spinner.Spinner{
		Frames: []string{
			"⡇ ",
			"⡇ ",
			"⣏ ",
			"⣉ ",
			" ⣉",
			" ⣹",
			" ⢸",
			" ⢸",
			" ⣹",
			" ⣉",
			"⣉ ",
			"⣏ ",
		},
		FPS: time.Second / 10,
	}

	stepStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("93")).Bold(true)
	pendingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	runningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("124"))
	canceledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))

	interactiveKeyStyle           = lipgloss.NewStyle().Background(lipgloss.Color("239")).Foreground(lipgloss.Color("152")).Padding(0, 1).Bold(true)
	interactiveTabStyle           = pendingStyle.Copy()
	interactiveActiveTabStyle     = runningStyle.Copy()
	interactiveFullScreenTabStyle = runningStyle.Copy().Background(lipgloss.Color("237"))
	interactiveHelpStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("249")).Faint(true)
	interactiveTaskStyle          = stepStyle.Copy()

	successIcon = successStyle.Render("✓")
	cancelIcon  = canceledStyle.Render("-")
	errorIcon   = errorStyle.Render("✗")
)

func HandleError(err error) {
	fmt.Println(errorStyle.Render("zwooc encountered an error"))
	fmt.Println(err)
	fmt.Println(runningStyle.Render("exiting zwooc"))
	os.Exit(1)
}

func PrintSuccess(msg string) {
	fmt.Printf(" %s %s\n", successStyle.Render("✓"), msg)
}
