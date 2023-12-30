package ui

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type notifyWriter struct {
	buf     *bytes.Buffer
	updates chan ContentUpdateMsg
}

var _ io.Writer = (*notifyWriter)(nil)

func NewNotifyWriter() *notifyWriter {
	return &notifyWriter{
		buf:     bytes.NewBuffer(nil),
		updates: make(chan ContentUpdateMsg, 8),
	}
}

func (w *notifyWriter) Write(p []byte) (n int, err error) {
	n, err = w.buf.Write(p)
	w.updates <- ContentUpdateMsg(w.buf.String())
	return n, err
}

func (w *notifyWriter) String() string {
	return w.buf.String()
}

func (w *notifyWriter) Close() error {
	close(w.updates)
	return nil
}

type ActiveTask struct {
	name   string
	writer *notifyWriter
	err    error
}

type Model struct {
	ready       bool
	updateCount int

	currentRunner tasks.TaskRunner

	activeTasks []ActiveTask
	logsView    viewport.Model

	// the $pre step of newly scheduled tasklists
	taskQueue []config.TaskList
	// merged tasklists of $post step of running tasks
	scheduledPost config.TaskList
}

type ContentUpdateMsg string

// NewInteractiveRunner creates a new interactive runner for long running tasks
func NewInteractiveRunner(list config.TaskList, opts ViewOptions, conf config.Config) error {
	_, main, _ := list.Split()
	out := NewNotifyWriter()
	for _, t := range main.Tasks {
		t.Pipe(out)
	}

	runner := tasks.NewRunner(main.Name, main.Tasks, opts.MaxConcurrency)
	go runner.Run()

	p := tea.NewProgram(
		&Model{
			writer: out,
		},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
	return nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.listenToUpdates, tea.EnterAltScreen)
}

func (m *Model) listenToUpdates() tea.Msg {
	return <-m.writer.updates
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case ContentUpdateMsg:
		m.logsView.SetContent(string(msg))
		m.logsView.GotoBottom()
		m.updateCount++
		cmds = append(cmds, m.listenToUpdates)
	case tea.WindowSizeMsg:
		// headerHeight := lipgloss.Height(m.headerView())
		// footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := 1 // headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.logsView = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.logsView.YPosition = 1
			m.logsView.HighPerformanceRendering = false // useHighPerformanceRenderer
			m.logsView.SetContent(m.writer.String())
			m.ready = true

			// // This is only necessary for high performance rendering, which in
			// // most cases you won't need.
			// //
			// // Render the viewport one line below the header.
			// m.logsView.YPosition = headerHeight + 1
		} else {
			m.logsView.Width = msg.Width
			m.logsView.Height = msg.Height - verticalMarginHeight
		}

		// if useHighPerformanceRenderer {
		// 	// Render (or re-render) the whole viewport. Necessary both to
		// 	// initialize the viewport and when the window is resized.
		// 	//
		// 	// This is needed for high-performance rendering only.
		// 	cmds = append(cmds, viewport.Sync(m.viewport))
		// }
	}

	// Handle keyboard and mouse events in the viewport
	m.logsView, cmd = m.logsView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s %d %d\n%s", "zwooc...", len(m.writer.String()), m.updateCount, m.logsView.View())
}
