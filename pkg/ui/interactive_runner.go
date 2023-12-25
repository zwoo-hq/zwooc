package ui

import (
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type taskStatus struct {
	name    string
	status  int
	spinner spinner.Model
}

type model struct {
	currentIndex  int
	tasks         config.TaskList
	tasksState    []taskStatus
	currentState  tasks.RunnerStatus
	currentRunner *tasks.TaskRunner
}

type updateMsg tasks.RunnerStatus
type stageFinishedMsg int
type errorMsg struct{ error }

func newInteractiveRunner(tasks config.TaskList) error {
	model := model{
		tasks:        tasks,
		currentIndex: 0,
	}

	execStart := time.Now()
	p := tea.NewProgram(&model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	execEnd := time.Now()
	fmt.Printf(" %s %s completed successfully in %s\n", successStyle.Render("✓"), tasks.Name, execEnd.Sub(execStart))
	return nil
}

func (m *model) Init() tea.Cmd {
	m.initStage(0)
	return tea.Batch(m.startStage, m.listenToUpdates, tea.EnterAltScreen)
}

func (m *model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
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
	case updateMsg:
		m.currentState = tasks.RunnerStatus(msg)
		m.convertRunnerState()

		cmds := []tea.Cmd{m.listenToUpdates}
		for _, task := range m.tasksState {
			if task.hasSpinner() {
				cmds = append(cmds, task.spinner.Tick)
			}
		}

		return m, tea.Batch(cmds...)
	case stageFinishedMsg:
		stage := int(msg)
		if stage+1 >= len(m.tasks.Steps) {
			return m, tea.Quit
		}
		m.initStage(stage + 1)
		return m, m.startStage
	case errorMsg:
		HandleError(msg.error)
		return m, tea.Quit
	}

	return m, nil
}

func (m *model) initStage(stage int) {
	m.currentIndex = stage
	m.currentRunner = tasks.NewRunner(m.tasks.Steps[stage].Name, m.tasks.Steps[stage].Tasks, m.tasks.Steps[stage].RunParallel)
	m.tasksState = []taskStatus{}
}

func (m *model) startStage() tea.Msg {
	err := m.currentRunner.Run()
	if err != nil {
		return errorMsg{err}
	}
	return stageFinishedMsg(m.currentIndex)
}

func (m *model) listenToUpdates() tea.Msg {
	return updateMsg(<-m.currentRunner.Updates())
}

func (m *model) View() (s string) {
	s += fmt.Sprintf("zwooc running: %s | %s (%d/%d)\n", m.tasks.Name, m.tasks.Steps[m.currentIndex].Name, m.currentIndex+1, len(m.tasks.Steps))
	s += "\n"

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

func (m *model) convertRunnerState() {
	if len(m.tasksState) == 0 {
		t := []taskStatus{}
		for key := range m.currentState {
			// set status to 0 to enforce a status update on first load
			t = append(t, taskStatus{name: key, status: 0})
		}
		sort.Slice(t, func(i, j int) bool {
			return t[i].name < t[j].name
		})
		m.tasksState = t
	}

	for i := 0; i < len(m.tasksState); i++ {
		status := &m.tasksState[i]
		newState, found := m.currentState[status.name]
		if !found {
			// stage changed
			m.tasksState = []taskStatus{}
			m.convertRunnerState()
			return
		}
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

func (t taskStatus) hasSpinner() bool {
	return t.status == tasks.StatusPending || t.status == tasks.StatusRunning
}
