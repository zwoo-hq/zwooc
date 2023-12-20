package ui

import (
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func init() {
	Logger = log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
		TimeFormat:      time.Kitchen,
		Prefix:          "zwooc",
		Level:           log.DebugLevel,
	})
	Logger.SetStyles(&log.Styles{
		Levels: map[log.Level]lipgloss.Style{
			log.DebugLevel: lipgloss.NewStyle().
				SetString("DEBUG").
				Padding(0, 1, 0, 1).
				Bold(true).
				Background(lipgloss.Color("252")).
				Foreground(lipgloss.Color("0")),
			log.InfoLevel: lipgloss.NewStyle().
				SetString("INFO").
				Padding(0, 1, 0, 2).
				Bold(true).
				Background(lipgloss.Color("081")).
				Foreground(lipgloss.Color("0")),
			log.WarnLevel: lipgloss.NewStyle().
				SetString("WARN").
				Padding(0, 1, 0, 2).
				Bold(true).
				Background(lipgloss.Color("220")).
				Foreground(lipgloss.Color("0")),
			log.ErrorLevel: lipgloss.NewStyle().
				SetString("ERROR").
				Padding(0, 1, 0, 1).
				Bold(true).
				Background(lipgloss.Color("196")).
				Foreground(lipgloss.Color("0")),
			log.FatalLevel: lipgloss.NewStyle().
				SetString("FATAL").
				Padding(0, 1, 0, 1).
				Bold(true).
				Background(lipgloss.Color("134")).
				Foreground(lipgloss.Color("0")),
		},
	})
}

func HandleError(err error) {
	Logger.Error(err.Error())
	Logger.Fatal("exiting zwooc")
}
