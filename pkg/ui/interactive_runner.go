package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type PreTaskStatus struct {
	name   string
	status int
	out    *tasks.CommandCapturer
}

type ActiveTask struct {
	name   string
	writer *notifyWriter
}

type ScheduledTask struct {
	preStage  tasks.TaskList
	mainTasks tasks.ExecutionStep
	postStage tasks.TaskList
}

type Model struct {
	err               error
	wasCanceled       bool
	wasCancelCanceled bool
	opts              ViewOptions
	scheduledTasks    []ScheduledTask

	preTasks         []PreTaskStatus
	preError         error
	preCurrentStage  int
	preCurrentList   tasks.TaskList
	preCurrentRunner *tasks.TaskRunner

	viewportReady bool
	activeTasks   []ActiveTask
	activeNotify  *notifyWriter
	taskToShow    string
	scheduler     *tasks.Scheduler
	logsView      viewport.Model

	scheduledPost     map[string]tasks.TaskList
	postError         error
	postTasks         []PreTaskStatus
	postCurrentStage  int
	postCurrentList   tasks.TaskList
	postCurrentRunner *tasks.TaskRunner
}

type ContentUpdateMsg string                // fired when the current logs content changes
type PreRunnerUpdateMsg tasks.RunnerStatus  //fired when a $pre tasks updates
type PostRunnerUpdateMsg tasks.RunnerStatus //fired when a $post tasks updates
type ScheduledStageFinishedMsg int          // fired whe a pre action of a scheduled task finished
type PostStageFinishedMsg int               // fired whe a post action of a scheduled task finished
type ScheduledErroredMsg struct{ error }    // fired when a scheduled task errored
type PostErroredMsg struct{ error }         // fired when a post task errored

// NewInteractiveRunner creates a new interactive runner for long running tasks
func NewInteractiveRunner(list tasks.TaskList, opts ViewOptions, conf config.Config) error {
	m := &Model{
		opts:           opts,
		scheduledTasks: []ScheduledTask{},
		activeTasks:    []ActiveTask{},
		scheduler:      tasks.NewScheduler(),
		scheduledPost:  make(map[string]tasks.TaskList),
	}

	m.schedule(list)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func (m *Model) Init() tea.Cmd {
	hasScheduledStage := m.prepareNextScheduled()
	if hasScheduledStage {
		return tea.Batch(tea.EnterAltScreen, m.startScheduledStage, m.listenToPreRunner)
	}
	if m.activeNotify != nil {
		return tea.Batch(tea.EnterAltScreen, m.listenToWriterUpdates)
	}
	return tea.Batch(tea.EnterAltScreen)
}

func (m *Model) schedule(t tasks.TaskList) {
	pre, main, post := t.Split()
	scheduled := ScheduledTask{
		preStage:  pre,
		mainTasks: main,
		postStage: post,
	}
	m.scheduledTasks = append(m.scheduledTasks, scheduled)
}

func (m *Model) prepareNextScheduled() bool {
	if len(m.scheduledTasks) == 0 {
		m.preCurrentList = tasks.TaskList{}
		m.preCurrentStage = 0
		m.preTasks = []PreTaskStatus{}
		m.preCurrentRunner = nil
		return false
	}

	scheduled := m.scheduledTasks[0]
	if scheduled.preStage.IsEmpty() {
		m.transitionCurrentScheduledIntoActive()
		m.scheduledTasks = m.scheduledTasks[1:]
		return m.prepareNextScheduled()
	}
	m.preCurrentStage = 0
	m.preCurrentList = scheduled.preStage
	m.initScheduledStage(0)
	return true
}

func (m *Model) initScheduledStage(stage int) {
	t := []PreTaskStatus{}
	for _, task := range m.preCurrentList.Steps[stage].Tasks {
		// set status to 0 to enforce a status update on first load
		cap := tasks.NewCapturer()
		task.Pipe(cap)
		t = append(t, PreTaskStatus{name: task.Name(), status: 0, out: cap})
	}
	sort.Slice(t, func(i, j int) bool {
		return t[i].name < t[j].name
	})

	m.preCurrentStage = stage
	m.preTasks = t
	m.preCurrentRunner = tasks.NewRunner(m.preCurrentList.Steps[stage].Name, m.preCurrentList.Steps[stage].Tasks, m.opts.MaxConcurrency)
}

func (m *Model) initPostStage(stage int) {
	t := []PreTaskStatus{}
	for _, task := range m.postCurrentList.Steps[stage].Tasks {
		// set status to 0 to enforce a status update on first load
		cap := tasks.NewCapturer()
		task.Pipe(cap)
		t = append(t, PreTaskStatus{name: task.Name(), status: 0, out: cap})
	}
	sort.Slice(t, func(i, j int) bool {
		return t[i].name < t[j].name
	})

	m.postCurrentStage = stage
	m.postTasks = t
	m.postCurrentRunner = tasks.NewRunner(m.postCurrentList.Steps[stage].Name, m.postCurrentList.Steps[stage].Tasks, m.opts.MaxConcurrency)
}

func (m *Model) listenToPreRunner() tea.Msg {
	return PreRunnerUpdateMsg(<-m.preCurrentRunner.Updates())
}

func (m *Model) listenToPostRunner() tea.Msg {
	return PostRunnerUpdateMsg(<-m.postCurrentRunner.Updates())
}

func (m *Model) startScheduledStage() tea.Msg {
	err := m.preCurrentRunner.Run()
	if err != nil {
		return ScheduledErroredMsg{err}
	}
	return ScheduledStageFinishedMsg(m.preCurrentStage)
}

func (m *Model) startPostStage() tea.Msg {
	err := m.postCurrentRunner.Run()
	if err != nil {
		return PostErroredMsg{err}
	}
	return PostStageFinishedMsg(m.preCurrentStage)
}

func (m *Model) transitionCurrentScheduledIntoActive() {
	current := m.scheduledTasks[0]
	if current.postStage.IsEmpty() && len(current.mainTasks.Tasks) == 0 {
		// this was a tasklist without a long running task - so it is already finished
		return
	}

	for _, task := range current.mainTasks.Tasks {
		notify := NewNotifyWriter()
		task.Pipe(notify)
		m.activeTasks = append(m.activeTasks, ActiveTask{name: task.Name(), writer: notify})
		m.scheduler.Schedule(task)
		if m.taskToShow == "" {
			// this is the first long running task
			m.taskToShow = current.mainTasks.Name
			m.activeNotify = notify
		}
	}

	m.scheduledPost[current.mainTasks.Name] = current.postStage
}

func (m *Model) listenToWriterUpdates() tea.Msg {
	return ContentUpdateMsg(<-m.activeNotify.updates)
}

func (m *Model) cancelAllRunning() tea.Msg {
	if m.wasCanceled {
		m.wasCancelCanceled = true
		return func() {}
	}

	m.wasCanceled = true
	m.err = m.scheduler.Cancel()
	if m.preCurrentRunner != nil {
		m.preError = m.preCurrentRunner.Cancel()
	}

	// reset active state
	m.activeTasks = []ActiveTask{}
	if m.activeNotify != nil {
		close(m.activeNotify.updates)
		m.activeNotify = nil
	}

	// start executing post tasks
	list := tasks.TaskList{
		Name: "cleanup",
	}
	for _, tasks := range m.scheduledPost {
		list.MergePostAligned(tasks)
	}
	list.RemoveEmptyStagesAndTasks()

	if list.IsEmpty() {
		return tea.Quit()
	}

	m.postCurrentStage = 0
	m.postCurrentList = list
	m.initPostStage(0)
	return m.startPostStage()
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
			return m, m.cancelAllRunning
		}
	case PreRunnerUpdateMsg:
		m.convertPreRunnerState(tasks.RunnerStatus(msg))
		if m.preCurrentRunner != nil {
			cmds = append(cmds, m.listenToPreRunner)
		}

	case PostRunnerUpdateMsg:
		m.convertPostRunnerState(tasks.RunnerStatus(msg))
		if m.postCurrentRunner != nil {
			cmds = append(cmds, m.listenToPostRunner)
		}

	case ScheduledStageFinishedMsg:
		stage := int(msg)
		if stage+1 >= len(m.preCurrentList.Steps) || m.wasCanceled {
			if !m.wasCanceled {
				m.transitionCurrentScheduledIntoActive()
			}
			m.scheduledTasks = m.scheduledTasks[1:]
			hasNext := m.prepareNextScheduled()
			if hasNext && !m.wasCanceled {
				cmds = append(cmds, m.startScheduledStage, m.listenToPreRunner)
			}
			if m.activeNotify != nil {
				cmds = append(cmds, m.listenToWriterUpdates)
			}
		} else if !m.wasCanceled {
			m.initScheduledStage(stage + 1)
			cmds = append(cmds, m.startScheduledStage, m.listenToPreRunner)
		}

	case PostStageFinishedMsg:
		stage := int(msg)
		if stage+1 >= len(m.postCurrentList.Steps) || m.wasCancelCanceled {
			return m, tea.Quit
		} else if !m.wasCancelCanceled {
			m.initPostStage(stage + 1)
			cmds = append(cmds, m.startPostStage, m.listenToPostRunner)
		}

	case ScheduledErroredMsg:
		// TODO: what to do here?
		m.preError = msg.error

	case PostErroredMsg:
		// TODO: what to do here?
		m.postError = msg.error

	case ContentUpdateMsg:
		m.logsView.SetContent(string(msg))
		m.logsView.GotoBottom()
		if m.activeNotify != nil {
			cmds = append(cmds, m.listenToWriterUpdates)
		}

	case tea.WindowSizeMsg:
		// headerHeight := lipgloss.Height(m.headerView())
		// footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := 8 // headerHeight + footerHeight

		if !m.viewportReady {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.logsView = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.logsView.YPosition = 8
			m.logsView.HighPerformanceRendering = false // useHighPerformanceRenderer
			// m.logsView.SetContent(m.writer.String())
			m.viewportReady = true
			m.logsView.SetContent("== empty ==")

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

func (m *Model) convertPreRunnerState(state tasks.RunnerStatus) {
	for i := 0; i < len(m.preTasks); i++ {
		status := &m.preTasks[i]
		newState := state[status.name]
		status.status = newState
	}
}

func (m *Model) convertPostRunnerState(state tasks.RunnerStatus) {
	for i := 0; i < len(m.postTasks); i++ {
		status := &m.postTasks[i]
		newState := state[status.name]
		status.status = newState
	}
}

func (m *Model) View() (s string) {
	header := fmt.Sprintf("zwooc running in interactive mode (%d scheduled tasks)\n", len(m.scheduledTasks))

	var currentTasks string
	if len(m.scheduledTasks) > 0 {
		currentlyRunning := []string{}
		for _, task := range m.preTasks {
			if task.status == tasks.StatusRunning {
				currentlyRunning = append(currentlyRunning, task.name)
			}
		}
		currentTasks = fmt.Sprintf("preparing %s [] running (%s)", m.scheduledTasks[0].mainTasks.Name, strings.Join(currentlyRunning, ", "))
	} else {
		currentTasks = "There are no tasks scheduled"
	}

	var postTasks string
	if !m.postCurrentList.IsEmpty() {
		currentlyRunning := []string{}
		for _, task := range m.postTasks {
			if task.status == tasks.StatusRunning {
				currentlyRunning = append(currentlyRunning, task.name)
			}
		}
		postTasks = fmt.Sprintf("shutting down %s [] running (%s)", m.postCurrentList.Name, strings.Join(currentlyRunning, ", "))
	} else {
		postTasks = "There are no tasks shutting down"
	}

	s += header
	s += "\n"
	s += currentTasks
	s += "\n\n"
	s += postTasks
	s += "\n\n"

	if !m.viewportReady {
		s += "Initializing..."
	} else if m.wasCanceled {
		s += "Shutting down..."
	} else {
		s += m.logsView.View()
	}
	return
}
