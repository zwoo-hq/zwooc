package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type InteractiveTaskStatus struct {
	name    string
	status  int
	out     *tasks.CommandCapturer
	spinner spinner.Model
}

type StatusModel struct {
	currentIndex  int
	tasks         tasks.TaskList
	tasksState    []InteractiveTaskStatus
	currentState  tasks.RunnerStatus
	currentRunner *tasks.TaskRunner
	opts          ViewOptions
	currentError  error
	wasCanceled   bool
	clear         bool
}

type StatusUpdateMsg tasks.RunnerStatus
type StatusStageFinishedMsg int
type StatusErrorMsg struct{ error }

func NewStatusView(list tasks.TaskList, opts ViewOptions) error {
	model := StatusModel{
		tasks:        list,
		currentIndex: 0,
		opts:         opts,
	}

	execStart := time.Now()
	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		return err
	}

	execEnd := time.Now()
	if model.currentError != nil {
		for _, status := range model.tasksState {
			if status.status == tasks.StatusError {
				fmt.Printf(" %s %s failed\n", errorStyle.Render("✗"), status.name)
				fmt.Printf(" %s error: %s\n", errorStyle.Render("✗"), model.currentError)
				fmt.Printf(" %s stdout:\n", errorStyle.Render("✗"))
				wrapper := canceledStyle.Render("===")
				parts := strings.Split(wrapper, "===")
				fmt.Printf(parts[0])
				fmt.Println(strings.TrimSpace(status.out.String()))
				fmt.Printf(parts[1])
				os.Exit(1)
			}
		}
		return nil
	}
	if model.wasCanceled {
		fmt.Printf("  %s %s canceled - stopping execution\n", canceledStyle.Render("-"), model.currentRunner.Name())
		return nil
	}
	fmt.Printf(" %s %s completed successfully in %s\n", successStyle.Render("✓"), list.Name, execEnd.Sub(execStart))
	return nil
}

func (m *StatusModel) Init() tea.Cmd {
	m.initStage(0)
	return tea.Batch(m.startStage, m.listenToUpdates)
}

func (m *StatusModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.wasCanceled = true
			m.currentRunner.Cancel()
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		for i := 0; i < len(m.tasksState); i++ {
			task := &m.tasksState[i]
			if task.spinner.ID() == msg.ID {
				task.spinner, cmd = task.spinner.Update(msg)
			}
		}
		return m, cmd
	case StatusUpdateMsg:
		m.currentState = tasks.RunnerStatus(msg)
		m.convertRunnerState()

		cmds := []tea.Cmd{m.listenToUpdates}
		for _, task := range m.tasksState {
			if task.hasSpinner() {
				cmds = append(cmds, task.spinner.Tick)
			}
		}

		return m, tea.Batch(cmds...)
	case StatusStageFinishedMsg:
		stage := int(msg)
		if stage+1 >= len(m.tasks.Steps) || m.wasCanceled {
			m.clear = true
			return m, tea.Quit
		}
		m.initStage(stage + 1)
		return m, m.startStage
	case StatusErrorMsg:
		m.currentError = msg.error
		m.clear = true
		return m, tea.Quit
	}

	return m, nil
}

func (m *StatusModel) initStage(stage int) {

	t := []InteractiveTaskStatus{}
	for _, task := range m.tasks.Steps[stage].Tasks {
		// set status to 0 to enforce a status update on first load
		cap := tasks.NewCapturer()
		task.Pipe(cap)
		t = append(t, InteractiveTaskStatus{name: task.Name(), status: 0, out: cap})
	}
	sort.Slice(t, func(i, j int) bool {
		return t[i].name < t[j].name
	})

	m.currentIndex = stage
	m.tasksState = t
	m.currentRunner = tasks.NewRunner(m.tasks.Steps[stage].Name, m.tasks.Steps[stage].Tasks, m.opts.MaxConcurrency)
}

func (m *StatusModel) startStage() tea.Msg {
	err := m.currentRunner.Run()
	if err != nil {
		return StatusErrorMsg{err}
	}
	return StatusStageFinishedMsg(m.currentIndex)
}

func (m *StatusModel) listenToUpdates() tea.Msg {
	return StatusUpdateMsg(<-m.currentRunner.Updates())
}

func (m *StatusModel) View() (s string) {
	if m.clear {
		return ""
	}

	s += fmt.Sprintf("zwooc running: %s | %s (%d/%d)\n\n", m.tasks.Name, m.tasks.Steps[m.currentIndex].Name, m.currentIndex+1, len(m.tasks.Steps))

	for _, task := range m.tasksState {
		if task.hasSpinner() {
			s += fmt.Sprintf(" %s %s: %s\n", task.spinner.View(), task.name, convertState(task.status))
		} else if task.status == tasks.StatusDone {
			s += fmt.Sprintf(" %s %s: %s\n", successStyle.Padding(0, 2).Render("✓"), task.name, convertState(task.status))
		} else if task.status == tasks.StatusError {
			s += fmt.Sprintf(" %s %s: %s\n", errorStyle.Padding(0, 2).Render("✗"), task.name, convertState(task.status))
		} else if task.status == tasks.StatusCanceled {
			s += fmt.Sprintf(" %s %s: %s\n", canceledStyle.Padding(0, 2).Render("-"), task.name, convertState(task.status))
		}
	}
	return
}

func convertState(state int) string {
	switch state {
	case tasks.StatusPending:
		return "pending"
	case tasks.StatusRunning:
		return "running"
	case tasks.StatusDone:
		return "done"
	case tasks.StatusError:
		return "error"
	case tasks.StatusCanceled:
		return "canceled"
	}
	return "unknown"
}

func (m *StatusModel) convertRunnerState() {
	for i := 0; i < len(m.tasksState); i++ {
		status := &m.tasksState[i]
		newState := m.currentState[status.name]
		if newState != status.status {
			status.status = newState
			status.spinner = spinner.New()
			switch newState {
			case tasks.StatusPending:
				status.spinner.Spinner = pendingSpinner
				status.spinner.Style = pendingStyle
			case tasks.StatusRunning:
				status.spinner.Spinner = runningSpinner
				status.spinner.Style = runningStyle
			}
		}
	}
}

func (t InteractiveTaskStatus) hasSpinner() bool {
	return t.status == tasks.StatusPending || t.status == tasks.StatusRunning
}
