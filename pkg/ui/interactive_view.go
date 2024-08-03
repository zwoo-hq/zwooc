package ui

import (
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zwoo-hq/zwooc/pkg/helper"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
	"github.com/zwoo-hq/zwooc/pkg/ui/textinput"
)

type activeViewType int

const (
	viewDefault activeViewType = iota
	viewHelp
	viewFullScreen
	viewAddTask
)

type taskUpdateMsg StatusUpdate
type runnerDoneMsg struct{ error }
type contentUpdateMsg struct {
	tabId   int
	content string
}

type interactiveTab struct {
	name     string
	writer   *tasks.NotifyWriter
	showLogs bool
	task     *tasks.TaskTreeNode
}

type interactiveView struct {
	tasks    tasks.Collection
	opts     ViewOptions
	provider *SchedulerStatusProvider

	outputs          map[string]*tasks.CommandCapturer
	status           map[string]TaskStatus
	aggregatedStatus map[string]TaskStatus

	viewportReady bool
	activeIndex   int
	tabs          []interactiveTab
	logsView      viewport.Model
	treeView      *treeProgressView

	input textinput.Model

	activeView   activeViewType
	windowWidth  int
	windowHeight int

	wasCanceled       bool
	wasCancelCanceled bool
	err               error
	clear             bool
}

func newInteractiveView(forest tasks.Collection, provider *SchedulerStatusProvider, opts ViewOptions) error {
	m := interactiveView{
		tasks:    forest,
		opts:     opts,
		provider: provider,

		status:           map[string]TaskStatus{},
		aggregatedStatus: map[string]TaskStatus{},
		outputs:          map[string]*tasks.CommandCapturer{},
		treeView: &treeProgressView{
			opts:    opts,
			spinner: map[TaskStatus]spinner.Model{},
		},

		tabs:        []interactiveTab{},
		activeIndex: -1,
		activeView:  viewDefault,

		input: textinput.New(),
	}

	m.treeView.status = m.status
	m.treeView.aggregatedStatus = m.aggregatedStatus

	m.input.Placeholder = "Enter a task key"
	m.input.Cursor.Style = interactiveActiveTabStyle
	m.input.Width = 30
	m.input.ShowSuggestions = true
	m.input.SetSuggestions([]string{"test", "test2", "test3"})

	execStart := time.Now()
	p := tea.NewProgram(&m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		return err
	}
	execEnd := time.Now()

	var failedError *tasks.MultiTaskError
	if errors.As(m.err, &failedError) {
		// handle runner error
		for nodeId, err := range failedError.Errors {
			fmt.Printf("%s %s failed: %s\n", errorIcon, nodeId, err)
		}
		fmt.Printf("%s %s %s failed after %s\n", zwoocBranding, errorIcon, forest.GetName(), execEnd.Sub(execStart))
	} else if m.wasCancelCanceled || errors.Is(m.err, tasks.ErrCancelled) {
		fmt.Printf("%s %s %s canceled after %s\n", zwoocBranding, cancelIcon, forest.GetName(), execEnd.Sub(execStart))
	} else {
		fmt.Printf("%s %s %s completed after  %s\n", zwoocBranding, successIcon, forest.GetName(), execEnd.Sub(execStart))
	}
	return nil
}

func (m *interactiveView) setupDefaultStatus() {
	for _, tree := range m.tasks {
		tree.Iterate(func(node *tasks.TaskTreeNode) {
			// set default status
			m.status[node.NodeID()] = StatusPending
			m.aggregatedStatus[node.NodeID()] = StatusPending
			// capture the output of each task
			cap := tasks.NewCapturer()
			m.outputs[node.NodeID()] = cap
			node.Main.Pipe(cap)
		})

		writer := tasks.NewNotifyWriter()
		tree.Main.Pipe(writer)
		m.tabs = append(m.tabs, interactiveTab{
			name:     tree.Name,
			writer:   writer,
			showLogs: false,
			task:     tree,
		})
	}

	if len(m.tabs) > 0 {
		m.activeIndex = 0
	}
}

func (m *interactiveView) Init() tea.Cmd {
	tea.SetWindowTitle("zwooc")

	m.setupDefaultStatus()

	return tea.Batch(m.listenToUpdates, m.start, m.treeView.setupSpinners())
}

func (m *interactiveView) updateProgress(update taskUpdateMsg) {
	m.status[update.NodeID] = update.Status
	m.aggregatedStatus[update.NodeID] = update.AggregatedStatus
	if update.Parent != nil {
		m.updateProgress(taskUpdateMsg(*update.Parent))
	}
}

func (m *interactiveView) listenToUpdates() tea.Msg {
	return taskUpdateMsg(<-m.provider.status)
}

func (m *interactiveView) start() tea.Msg {
	m.provider.Start()
	return runnerDoneMsg{<-m.provider.done}
}

func (m *interactiveView) listenToWriterUpdates() tea.Msg {
	currentIdx := m.activeIndex
	if currentIdx < 0 || currentIdx >= len(m.tabs) || !m.tabs[currentIdx].showLogs {
		return nil
	}

	return contentUpdateMsg{
		tabId:   currentIdx,
		content: <-m.tabs[currentIdx].writer.Updates,
	}
}

func (m *interactiveView) updateCurrentLogsView() tea.Msg {
	if m.activeIndex < 0 || m.activeIndex >= len(m.tabs) {
		return nil
	}

	if m.tabs[m.activeIndex].showLogs {
		return contentUpdateMsg{
			tabId:   m.activeIndex,
			content: m.tabs[m.activeIndex].writer.String(),
		}
	}

	return contentUpdateMsg{
		tabId:   m.activeIndex,
		content: m.treeView.printNode(m.tabs[m.activeIndex].task, "", true),
	}
}

func (m *interactiveView) handleCancel() {
	if m.wasCanceled && !m.wasCancelCanceled {
		m.wasCancelCanceled = true
		m.provider.Cancel()
	} else {
		m.wasCanceled = true
		for i := range m.tabs {
			m.tabs[i].showLogs = false
		}
		m.provider.Cancel()
	}
}

func (m *interactiveView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.handleCancel()
		case "h":
			if m.activeView == viewHelp {
				m.activeView = viewDefault
				m.setLogsViewDefaultPosition()
			} else {
				m.activeView = viewHelp
			}
		case "f":
			if m.activeView == viewFullScreen {
				m.activeView = viewDefault
				m.setLogsViewDefaultPosition()
			} else {
				m.activeView = viewFullScreen
				m.setLogsViewFullScreenPosition()
			}
		case "a":
			if m.activeView == viewAddTask {
				m.activeView = viewDefault
				m.setLogsViewDefaultPosition()
			} else {
				m.activeView = viewAddTask
				m.input.Focus()
				cmds = append(cmds, textinput.Blink)
			}
		case "esc":
			m.activeView = viewDefault
			m.setLogsViewDefaultPosition()
		case "tab":
			if len(m.tabs) > 0 {
				m.activeIndex = (m.activeIndex + 1) % len(m.tabs)
				cmds = append(cmds, m.listenToWriterUpdates, m.updateCurrentLogsView)
			}
		case "shift+tab":
			if len(m.tabs) > 0 {
				m.activeIndex = (m.activeIndex - 1 + len(m.tabs)) % len(m.tabs)
				cmds = append(cmds, m.listenToWriterUpdates, m.updateCurrentLogsView)
			}
		}

	case spinner.TickMsg:
		_, cmd = m.treeView.Update(msg)
		return m, cmd

	case taskUpdateMsg:
		m.updateProgress(msg)

		for i, tab := range m.tabs {
			preNodes := helper.MapTo(tab.task.Pre, func(node *tasks.TaskTreeNode) TaskStatus {
				return m.aggregatedStatus[node.NodeID()]
			})
			m.tabs[i].showLogs = helper.All(preNodes, func(status TaskStatus) bool {
				return status == StatusDone
			})
		}

		return m, tea.Batch(m.listenToUpdates, m.updateCurrentLogsView)
	case runnerDoneMsg:
		m.err = msg.error
		m.clear = true
		return m, tea.Quit

	case contentUpdateMsg:
		// this is to ignore old (pending) updates from other tabs after the tab changed
		if msg.tabId == m.activeIndex {
			m.logsView.SetContent(string(msg.content))
			m.logsView.GotoBottom()
			if m.activeIndex >= 0 {
				cmds = append(cmds, m.listenToWriterUpdates)
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft && msg.Y > 4 && msg.Y < 8 && m.activeView == viewDefault {
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
			if m.activeView == viewFullScreen {
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

func (m *interactiveView) setLogsViewDefaultPosition() {
	m.logsView.Width = m.windowWidth
	m.logsView.Height = m.windowHeight - 9
}

func (m *interactiveView) setLogsViewFullScreenPosition() {
	m.logsView.Width = m.windowWidth
	m.logsView.Height = m.windowHeight - 1
}

func (m *interactiveView) View() (s string) {
	if m.clear {
		return ""
	}

	if m.activeView == viewHelp {
		return m.ViewHelp()
	}

	if m.activeView == viewFullScreen {
		return m.ViewFullScreen()
	}

	if m.activeView == viewAddTask {
		return m.ViewAddTask()
	}

	header := "zwooc running in interactive mode\n"

	var currentTasks string
	var postTasks string

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

func (m *interactiveView) ViewHelp() (s string) {
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

func (m *interactiveView) ViewFullScreen() (s string) {
	if m.activeIndex < 0 || len(m.tabs) == 0 {
		return "there is no active tab"
	}
	name := interactiveFullScreenTabStyle.Render(" " + m.tabs[m.activeIndex].name + " ")
	help := interactiveKeyStyle.Render("h") + interactiveHelpStyle.Render(" • show help")
	fs := interactiveKeyStyle.Render("f") + interactiveHelpStyle.Render(" • toggle fullscreen")

	start := fmt.Sprintf("╾─┤%s├", name)
	end := fmt.Sprintf(" %s │ %s ", fs, help)
	middle := helper.Repeat("─", m.logsView.Width-lipgloss.Width(start)-lipgloss.Width(end)-2)

	s += start + middle + "─╼" + end + "\n"
	s += m.logsView.View()
	return
}

func (m *interactiveView) ViewAddTask() (s string) {
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

func (m *interactiveView) RenderTabs() string {
	tabsTop := "╭─"
	tabs := "│ "
	tabsBorder := "┵─"

	for i, task := range m.tabs {
		var currentName string
		if i == m.activeIndex {
			currentName = interactiveActiveTabStyle.Render(task.name)
		} else {
			currentName = interactiveTabStyle.Render(task.name)
		}
		tabs += currentName + " │ "
		tabsBorder += helper.Repeat("─", lipgloss.Width(currentName)) + "─┴─"
		tabsTop += helper.Repeat("─", lipgloss.Width(currentName)) + "─"
		if i == len(m.tabs)-1 {
			tabsTop += "╮"
		} else {
			tabsTop += "┬─"
		}
	}

	if len(m.tabs) == 0 {
		tabsTop = "╭───────────────────╮"
		tabs = "│ (no active tasks) │"
		tabsBorder = "┵───────────────────┴"
	}
	help := interactiveKeyStyle.Render("tab") + interactiveHelpStyle.Render(" • switch tab")
	tabsBorder += helper.Repeat("─", m.logsView.Width-3-lipgloss.Width(tabsBorder)-lipgloss.Width(help))
	tabsBorder += "┤ " + help

	return tabsTop + "\n" + tabs + "\n" + tabsBorder + "\n"
}

func (m *interactiveView) determineTabClicked(x int) int {
	var current = 0
	for i, task := range m.tabs {
		tabWidth := len(task.name) + 2
		if x > current && x < current+tabWidth+1 {
			return i
		}
		current += tabWidth + 1
	}

	return -1
}
