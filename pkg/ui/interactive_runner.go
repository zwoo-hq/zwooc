package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zwoo-hq/zwooc/pkg/config"
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/ui/textinput"
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

type ActiveView int

const (
	ViewDefault ActiveView = iota
	ViewHelp
	ViewFullScreen
	ViewAddTask
)

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
	preSpinner       spinner.Model

	viewportReady bool
	activeTasks   []ActiveTask
	activeIndex   int
	scheduler     *tasks.Scheduler
	logsView      viewport.Model

	scheduledPost     map[string]tasks.TaskList
	postError         error
	postTasks         []PreTaskStatus
	postCurrentStage  int
	postCurrentList   tasks.TaskList
	postCurrentRunner *tasks.TaskRunner
	postSpinner       spinner.Model

	input textinput.Model

	activeView   ActiveView
	windowWidth  int
	windowHeight int
}

// fired when the current logs content changes
type ContentUpdateMsg struct {
	tabId   int
	content string
}
type PreRunnerUpdateMsg tasks.RunnerStatus  // fired when a $pre tasks updates
type PostRunnerUpdateMsg tasks.RunnerStatus // fired when a $post tasks updates
type ScheduledStageFinishedMsg int          // fired when a pre action of a scheduled task finished
type PostStageFinishedMsg int               // fired when a post action of a scheduled task finished
type ScheduledErroredMsg struct{ error }    // fired when a scheduled task errored
type PostErroredMsg struct{ error }         // fired when a post task errored

// NewInteractiveRunner creates a new interactive runner for long running tasks
func NewInteractiveRunner(forest tasks.Collection, opts ViewOptions, conf config.Config) error {
	m := &Model{
		opts:           opts,
		scheduledTasks: []ScheduledTask{},
		activeTasks:    []ActiveTask{},
		scheduler:      tasks.NewScheduler(),
		scheduledPost:  make(map[string]tasks.TaskList),
		activeIndex:    -1,
		activeView:     ViewDefault,
		input:          textinput.New(),
		preSpinner:     spinner.New(),
		postSpinner:    spinner.New(),
	}

	m.preSpinner.Spinner = pendingTabSpinner
	m.preSpinner.Style = pendingStyle

	m.postSpinner.Spinner = shutdownTabSpinner
	m.postSpinner.Style = errorStyle

	m.input.Placeholder = "Enter a task key"
	m.input.Cursor.Style = interactiveActiveTabStyle
	m.input.Width = 30
	m.input.ShowSuggestions = true
	m.input.SetSuggestions([]string{"test", "test2", "test3"})

	for _, tree := range forest {
		list := tree.Flatten()
		m.schedule(list)
	}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func (m *Model) Init() tea.Cmd {
	tea.SetWindowTitle("zwooc")

	hasScheduledStage := m.prepareNextScheduled()
	if hasScheduledStage {
		return tea.Batch(m.startScheduledStage, m.listenToPreRunner, m.preSpinner.Tick, m.postSpinner.Tick)
	}
	if len(m.activeTasks) > 0 {
		m.activeIndex = 0
		return tea.Batch(m.listenToWriterUpdates, m.preSpinner.Tick, m.postSpinner.Tick)
	}
	return tea.Batch(m.preSpinner.Tick, m.postSpinner.Tick)
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
	if m.preCurrentRunner == nil {
		return func() {}
	}
	return PreRunnerUpdateMsg(<-m.preCurrentRunner.Updates())
}

func (m *Model) listenToPostRunner() tea.Msg {
	if m.postCurrentRunner == nil {
		return func() {}
	}
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
	return PostStageFinishedMsg(m.postCurrentStage)
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
		if m.activeIndex < 0 {
			// set this as current tab
			m.activeIndex = len(m.activeTasks) - 1
		}
	}

	m.scheduledPost[current.mainTasks.Name] = current.postStage
}

func (m *Model) listenToWriterUpdates() tea.Msg {
	currentId := m.activeIndex
	if currentId < 0 || currentId >= len(m.activeTasks) {
		return func() {}
	}
	return ContentUpdateMsg{
		tabId:   currentId,
		content: <-m.activeTasks[currentId].writer.updates,
	}
}

func (m *Model) updateCurrentLogsView() tea.Msg {
	return ContentUpdateMsg{
		tabId:   m.activeIndex,
		content: m.activeTasks[m.activeIndex].writer.String(),
	}
}

func (m *Model) cancelAllRunning() bool {
	if m.wasCanceled {
		m.wasCancelCanceled = true
		return false
	}

	m.wasCanceled = true
	m.err = m.scheduler.Cancel()
	if m.preCurrentRunner != nil {
		m.preError = m.preCurrentRunner.Cancel()
	}

	// reset active state
	m.activeTasks = []ActiveTask{}
	m.activeIndex = -1
	m.logsView.SetContent("Shutting down...\n")

	// start executing post tasks
	list := tasks.TaskList{
		Name: "cleanup",
	}
	for _, tasks := range m.scheduledPost {
		list.MergePostAligned(tasks)
	}
	list.RemoveEmptyStagesAndTasks()

	if list.IsEmpty() {
		return false
	}

	m.postCurrentStage = 0
	m.postCurrentList = list
	m.initPostStage(0)
	return true
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
			if m.cancelAllRunning() {
				return m, tea.Batch(m.startPostStage, m.listenToPostRunner)
			} else {
				return m, tea.Quit
			}
		case "h":
			if m.activeView == ViewHelp {
				m.activeView = ViewDefault
				m.setLogsViewDefaultPosition()
			} else {
				m.activeView = ViewHelp
			}
		case "f":
			if m.activeView == ViewFullScreen {
				m.activeView = ViewDefault
				m.setLogsViewDefaultPosition()
			} else {
				m.activeView = ViewFullScreen
				m.setLogsViewFullScreenPosition()
			}
		case "a":
			if m.activeView == ViewAddTask {
				m.activeView = ViewDefault
				m.setLogsViewDefaultPosition()
			} else {
				m.activeView = ViewAddTask
				m.input.Focus()
				cmds = append(cmds, textinput.Blink)
			}
		case "esc":
			m.activeView = ViewDefault
			m.setLogsViewDefaultPosition()
		case "tab":
			if len(m.activeTasks) > 0 {
				m.activeIndex = (m.activeIndex + 1) % len(m.activeTasks)
				cmds = append(cmds, m.listenToWriterUpdates, m.updateCurrentLogsView)
			}
		case "shift+tab":
			if len(m.activeTasks) > 0 {
				m.activeIndex = (m.activeIndex - 1 + len(m.activeTasks)) % len(m.activeTasks)
				cmds = append(cmds, m.listenToWriterUpdates, m.updateCurrentLogsView)
			}
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
			if m.activeIndex >= 0 {
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
		// this is to ignore old (pending) updates from other tabs after the tab changed
		if msg.tabId == m.activeIndex {
			m.logsView.SetContent(string(msg.content))
			m.logsView.GotoBottom()
			if m.activeIndex >= 0 {
				cmds = append(cmds, m.listenToWriterUpdates)
			}
		}

	case spinner.TickMsg:
		if m.preSpinner.ID() == msg.ID {
			m.preSpinner, cmd = m.preSpinner.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.postSpinner.ID() == msg.ID {
			m.postSpinner, cmd = m.postSpinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft && msg.Y > 4 && msg.Y < 8 && m.activeView == ViewDefault {
			clickedIdx := m.determineTabClicked(msg.X)
			if clickedIdx >= 0 {
				m.activeIndex = clickedIdx
				cmds = append(cmds, m.listenToWriterUpdates, m.updateCurrentLogsView)
			}
		}

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

		if !m.viewportReady {
			m.logsView = viewport.New(msg.Width, msg.Height)
			m.logsView.HighPerformanceRendering = false
			// m.logsView.YPosition = 10 (use only with high performance rendering)
			m.viewportReady = true
			m.logsView.SetContent("== empty ==")
			m.setLogsViewDefaultPosition()
		} else {
			if m.activeView == ViewFullScreen {
				m.setLogsViewFullScreenPosition()
			} else {
				m.setLogsViewDefaultPosition()
			}
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
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) setLogsViewDefaultPosition() {
	m.logsView.Width = m.windowWidth
	m.logsView.Height = m.windowHeight - 9
}

func (m *Model) setLogsViewFullScreenPosition() {
	m.logsView.Width = m.windowWidth
	m.logsView.Height = m.windowHeight - 1
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
	if m.activeView == ViewHelp {
		return m.ViewHelp()
	}

	if m.activeView == ViewFullScreen {
		return m.ViewFullScreen()
	}

	if m.activeView == ViewAddTask {
		return m.ViewAddTask()
	}

	header := fmt.Sprintf("zwooc running in interactive mode (%d scheduled tasks)\n", len(m.scheduledTasks))

	var currentTasks string
	if len(m.scheduledTasks) > 0 {
		currentlyRunning := []string{}
		for _, task := range m.preTasks {
			if task.status == tasks.StatusRunning {
				currentlyRunning = append(currentlyRunning, task.name)
			}
		}
		currentTasks = fmt.Sprintf("%s preparing %s running [%s]", m.preSpinner.View(), interactiveTaskStyle.Render(m.scheduledTasks[0].mainTasks.Name), strings.Join(currentlyRunning, ", "))
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
		postTasks = fmt.Sprintf("%s shutting down %s running [%s]", m.postSpinner.View(), interactiveTaskStyle.Render(m.postCurrentList.Name), strings.Join(currentlyRunning, ", "))
	} else {
		postTasks = "There are no tasks shutting down"
	}

	s += header
	s += "\n"
	s += currentTasks
	s += "\n\n"
	s += postTasks
	s += "\n"

	s += m.RenderTabs()

	if !m.viewportReady {
		s += "Initializing...\n"
	} else {
		s += m.logsView.View() + "\n"
	}
	help := interactiveKeyStyle.Render("h") + interactiveHelpStyle.Render(" • show help")
	s += fmt.Sprintf("╾%s┤ %s", helper.Repeat("─", m.logsView.Width-lipgloss.Width(help)-3), help)
	return
}

func (m *Model) ViewHelp() (s string) {
	s += "zwooc interactive runner - help\n\n"
	align := lipgloss.NewStyle().Width(12).Align(lipgloss.Right).MarginRight(1).MarginLeft(2)

	s += align.Render(interactiveKeyStyle.Render("q/ctrl+c")) + interactiveHelpStyle.Render(" quit the runner") + "\n\n"
	s += align.Render(interactiveKeyStyle.Render("h")) + interactiveHelpStyle.Render(" show/hide this help") + "\n\n"
	s += align.Render(interactiveKeyStyle.Render("f")) + interactiveHelpStyle.Render(" toggle full screen mode") + "\n\n"
	s += align.Render(interactiveKeyStyle.Render("esc")) + interactiveHelpStyle.Render(" close the alt (help) screen") + "\n\n"
	s += align.Render(interactiveKeyStyle.Render("tab")) + interactiveHelpStyle.Render(" switch to next tab") + "\n\n"
	s += align.Render(interactiveKeyStyle.Render("shift+tab")) + interactiveHelpStyle.Render(" switch to previous tab") + "\n\n"
	return
}

func (m *Model) ViewFullScreen() (s string) {
	if m.activeIndex < 0 || len(m.activeTasks) == 0 {
		return "there is no active tab"
	}
	name := interactiveFullScreenTabStyle.Render(" " + m.activeTasks[m.activeIndex].name + " ")
	help := interactiveKeyStyle.Render("h") + interactiveHelpStyle.Render(" • show help")
	fs := interactiveKeyStyle.Render("f") + interactiveHelpStyle.Render(" • toggle fullscreen")

	start := fmt.Sprintf("╾─┤%s├", name)
	end := fmt.Sprintf(" %s │ %s ", fs, help)
	middle := helper.Repeat("─", m.logsView.Width-lipgloss.Width(start)-lipgloss.Width(end)-2)

	s += start + middle + "─╼" + end + "\n"
	s += m.logsView.View()
	return
}

func (m *Model) ViewAddTask() (s string) {
	s += "zwooc interactive runner - add task\n\n"
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Width(m.windowWidth - 2)
	truncatedContent := lipgloss.NewStyle().MaxWidth(m.windowWidth - 2).Render(m.input.View())
	s += border.Render(truncatedContent)
	s += "\n"

	suggestions := m.input.MatchedSuggestions()
	if len(suggestions) == 0 {
		suggestions = m.input.AvailableSuggestions()
	}
	for _, suggestion := range suggestions {
		if m.input.CurrentSuggestion() == suggestion {
			s += "  ◦ " + interactiveActiveTabStyle.Render(suggestion) + "\n"
		} else {
			s += "  ◦ " + suggestion + "\n"
		}
	}
	return
}

func (m *Model) RenderTabs() string {
	tabsTop := "╭─"
	tabs := "│ "
	tabsBorder := "┵─"

	for i, task := range m.activeTasks {
		var currentName string
		if i == m.activeIndex {
			currentName = interactiveActiveTabStyle.Render(task.name)
		} else {
			currentName = interactiveTabStyle.Render(task.name)
		}
		tabs += currentName + " │ "
		tabsBorder += helper.Repeat("─", lipgloss.Width(currentName)) + "─┴─"
		tabsTop += helper.Repeat("─", lipgloss.Width(currentName)) + "─"
		if i == len(m.activeTasks)-1 {
			tabsTop += "╮"
		} else {
			tabsTop += "┬─"
		}
	}

	if len(m.activeTasks) == 0 {
		tabsTop = "╭───────────────────╮"
		tabs = "│ (no active tasks) │"
		tabsBorder = "┵───────────────────┴"
	}
	help := interactiveKeyStyle.Render("tab") + interactiveHelpStyle.Render(" • switch tab")
	tabsBorder += helper.Repeat("─", m.logsView.Width-3-lipgloss.Width(tabsBorder)-lipgloss.Width(help))
	tabsBorder += "┤ " + help

	return tabsTop + "\n" + tabs + "\n" + tabsBorder + "\n"
}

func (m *Model) determineTabClicked(x int) int {
	var current = 0
	for i, task := range m.activeTasks {
		tabWidth := len(task.name) + 2
		if x > current && x < current+tabWidth+1 {
			return i
		}
		current += tabWidth + 1
	}

	return -1
}
