package ui

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

	zwoocBranding = lipgloss.NewStyle().Foreground(lipgloss.CompleteColor{
		TrueColor: "#29A4E7",
		ANSI256:   "39",
		ANSI:      "14",
	}).Background(lipgloss.CompleteColor{
		TrueColor: "#181521",
		ANSI256:   "17",
		ANSI:      "0",
	}).Bold(true).Padding(0, 2).Render("⚡ zwooc")
	pendingStyle                  = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	runningStyle                  = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	successStyle                  = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	errorStyle                    = lipgloss.NewStyle().Foreground(lipgloss.Color("124"))
	canceledStyle                 = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	stepStyle                     = lipgloss.NewStyle().Foreground(lipgloss.Color("93")).Bold(true)
	graphHeaderStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("93")).Bold(true)
	graphMainStyle                = lipgloss.NewStyle().Foreground(lipgloss.Color("93"))
	graphPreStyle                 = lipgloss.NewStyle().Foreground(lipgloss.Color("117")).Italic(true)
	graphPostStyle                = lipgloss.NewStyle().Foreground(lipgloss.Color("130")).Italic(true)
	graphInfoStyle                = lipgloss.NewStyle().Foreground(lipgloss.Color("249")).Faint(true)
	interactiveKeyStyle           = lipgloss.NewStyle().Background(lipgloss.Color("239")).Foreground(lipgloss.Color("152")).Padding(0, 1).Bold(true)
	interactiveTabStyle           = pendingStyle.Copy()
	interactiveActiveTabStyle     = runningStyle.Copy()
	interactiveFullScreenTabStyle = runningStyle.Copy().Background(lipgloss.Color("237"))
	interactiveHelpStyle          = graphInfoStyle.Copy()
	interactiveTaskStyle          = stepStyle.Copy()
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
