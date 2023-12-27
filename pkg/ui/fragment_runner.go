package ui

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

func NewFragmentRunner(task tasks.Task, opts ViewOptions) {
	cmdOut := tasks.NewCapturer()
	task.Pipe(cmdOut)
	if !opts.QuiteMode {
		fmt.Printf("running %s\n", task.Name())
		task.Pipe(os.Stdout)
	}

	interrupt := make(chan os.Signal, 1)
	cancel := make(chan bool, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		for range interrupt {
			cancel <- true
			break
		}
	}()

	execStart := time.Now()
	err := task.Run(cancel)
	execEnd := time.Now()

	if err != nil {
		fmt.Printf(" %s %s failed\n", errorStyle.Render("✗"), task.Name())
		fmt.Printf(" %s error: %s\n", errorStyle.Render("✗"), err)
		fmt.Printf(" %s stdout:\n", errorStyle.Render("✗"))
		fmt.Println(canceledStyle.Render(strings.TrimSpace(cmdOut.String())))
		return
	}
	fmt.Printf(" %s %s completed successfully in %s\n", successStyle.Render("✓"), task.Name(), execEnd.Sub(execStart))
}
