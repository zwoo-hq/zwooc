package ui

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zwoo-hq/zwooc/pkg/model"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type TreeProgressView struct {
	tasks            tasks.Collection
	opts             ViewOptions
	outputs          map[string]*tasks.CommandCapturer
	status           map[string]TaskStatus
	aggregatedStatus map[string]TaskStatus
	provider         SimpleStatusProvider
	mu               sync.RWMutex
	wasCanceled      bool
	err              error
	clear            bool
	spinner          map[TaskStatus]spinner.Model
}

type TreeProgressUpdateMsg StatusUpdate
type TreeProgressDoneMsg struct{ error }

func NewTreeProgressView(forest tasks.Collection, status SimpleStatusProvider, opts ViewOptions) error {
	model := TreeProgressView{
		opts:             opts,
		tasks:            forest,
		provider:         status,
		status:           map[string]TaskStatus{},
		aggregatedStatus: map[string]TaskStatus{},
		outputs:          map[string]*tasks.CommandCapturer{},
		spinner:          map[TaskStatus]spinner.Model{},
	}

	model.setupDefaultStatus()
	model.setupInterruptHandler()

	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		return err
	}

	// TODO: done -display cancel or error or success

	fmt.Println("done!!")
	return nil
}

func (m *TreeProgressView) Init() tea.Cmd {
	return tea.Batch(m.listenToUpdates, m.start, m.setupSpinners())
}

func (m *TreeProgressView) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.provider.Cancel()
			m.mu.Lock()
			m.wasCanceled = true
			m.mu.Unlock()
			return m, nil
		}
	case spinner.TickMsg:
		var cmds []tea.Cmd
		for i, s := range m.spinner {
			newModel, cmd := s.Update(msg)
			m.spinner[i] = newModel
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case TreeProgressUpdateMsg:
		m.mu.Lock()
		m.updateProgress(msg)
		m.mu.Unlock()
		return m, m.listenToUpdates
	case TreeProgressDoneMsg:
		m.mu.Lock()
		m.err = msg.error
		m.clear = true
		m.mu.Unlock()
		return m, tea.Quit
	}

	return m, nil
}

func (m *TreeProgressView) updateProgress(update TreeProgressUpdateMsg) {
	m.status[update.NodeID] = update.Status
	m.aggregatedStatus[update.NodeID] = update.AggregatedStatus
	if update.Parent != nil {
		m.updateProgress(TreeProgressUpdateMsg(*update.Parent))
	}
}

func (m *TreeProgressView) listenToUpdates() tea.Msg {
	return TreeProgressUpdateMsg(<-m.provider.status)
}

func (m *TreeProgressView) start() tea.Msg {
	m.provider.Start()
	return TreeProgressDoneMsg{<-m.provider.done}
}

type X struct {
	N string
	S TaskStatus
}

func (m *TreeProgressView) View() (s string) {
	if m.clear {
		return
	}

	s += zwoocBranding
	s += "- executing " + m.tasks.GetName() + "\n"
	for i, tree := range m.tasks {
		s += m.printNode(tree, "", i == len(m.tasks)-1)
	}
	return
}

func (m *TreeProgressView) setupDefaultStatus() {
	for _, tree := range m.tasks {
		tree.Iterate(func(node *tasks.TaskTreeNode) {
			// set default status
			m.status[node.NodeID()] = StatusPending
			m.aggregatedStatus[node.NodeID()] = StatusPending
			// capture the output of each task
			cap := tasks.NewCapturer()
			m.outputs[node.NodeID()] = cap
			node.Main.Pipe(cap)
			// TODO: handle inline output
			// if opts.InlineOutput {
			// 	if opts.DisablePrefix {
			// 		node.Main.Pipe(tasks.NewPrefixer("│  ", os.Stdout))
			// 	} else {
			// 		node.Main.Pipe(tasks.NewPrefixer("│  "+node.Main.Name()+" ", os.Stdout))
			// 	}
			// }
		})
	}
}

func (m *TreeProgressView) setupInterruptHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		m.provider.Cancel()
		m.mu.Lock()
		m.wasCanceled = true
		m.mu.Unlock()
	}()
}

func (m *TreeProgressView) setupSpinners() tea.Cmd {
	pendingSpinner := spinner.New(spinner.WithSpinner(pendingTabSpinner), spinner.WithStyle(treePendingStyle))
	scheduledSpinner := spinner.New(spinner.WithSpinner(pendingTabSpinner), spinner.WithStyle(treeScheduledStyle))
	runningSpinner := spinner.New(spinner.WithSpinner(runningTabSpinner), spinner.WithStyle(treeRunningStyle))

	m.spinner[StatusPending] = pendingSpinner
	m.spinner[StatusScheduled] = scheduledSpinner
	m.spinner[StatusRunning] = runningSpinner
	m.spinner[StatusDone] = spinner.New(spinner.WithSpinner(spinner.Spinner{
		Frames: []string{"✓ "},
		FPS:    1,
	}), spinner.WithStyle(treeSuccessStyle))
	m.spinner[StatusError] = spinner.New(spinner.WithSpinner(spinner.Spinner{
		Frames: []string{"✗ "},
		FPS:    1,
	}), spinner.WithStyle(treeErrorStyle))
	m.spinner[StatusCanceled] = spinner.New(spinner.WithSpinner(spinner.Spinner{
		Frames: []string{"- "},
		FPS:    1,
	}), spinner.WithStyle(treeCanceledStyle))

	return tea.Batch(scheduledSpinner.Tick, runningSpinner.Tick, pendingSpinner.Tick)
}

func (m *TreeProgressView) printNode(node *tasks.TaskTreeNode, prefix string, isLast bool) (s string) {
	connector := "┬"
	status := m.aggregatedStatus[node.NodeID()]
	if node.IsLeaf() {
		connector = "─"
		status = m.status[node.NodeID()]
	}

	nodeStatus := m.spinner[status].View()
	if isLast {
		s += fmt.Sprintf("%s└%s%s %s\n", prefix, connector, node.Name, nodeStatus)
	} else {
		s += fmt.Sprintf("%s├%s%s %s\n", prefix, connector, node.Name, nodeStatus)
	}

	if node.IsLeaf() {
		return
	}

	descendantPrefix := "│"
	if isLast {
		descendantPrefix = " "
	}

	if len(node.Pre) > 0 {
		s += fmt.Sprintf("%s%s├┬%s\n", prefix, descendantPrefix, model.KeyPre)
		for i, child := range node.Pre {
			s += m.printNode(child, prefix+descendantPrefix+"│", i == len(node.Pre)-1)
		}
	}

	mainStatus := m.spinner[m.status[node.NodeID()]].View()
	if len(node.Post) > 0 {
		s += fmt.Sprintf("%s%s├─%s %s\n", prefix, descendantPrefix, node.Main.Name(), mainStatus)
	} else {
		s += fmt.Sprintf("%s%s└─%s %s\n", prefix, descendantPrefix, node.Main.Name(), mainStatus)
	}

	if len(node.Post) > 0 {
		s += fmt.Sprintf("%s%s└┬%s\n", prefix, descendantPrefix, model.KeyPost)
		for i, child := range node.Post {
			s += m.printNode(child, prefix+descendantPrefix+" ", i == len(node.Post)-1)
		}
	}

	return
}
