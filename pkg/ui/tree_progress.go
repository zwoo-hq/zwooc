package ui

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zwoo-hq/zwooc/pkg/tasks"
)

type TreeProgressView struct {
	tasks       tasks.Collection
	opts        ViewOptions
	outputs     map[string]*tasks.CommandCapturer
	status      map[string]TaskStatus
	provider    SimpleStatusProvider
	mu          sync.RWMutex
	wasCanceled bool
}

type TreeProgressUpdateMsg StatusUpdate

func NewTreeProgressView(forest tasks.Collection, status SimpleStatusProvider, opts ViewOptions) error {
	model := TreeProgressView{
		opts:     opts,
		tasks:    forest,
		provider: status,
		status:   map[string]TaskStatus{},
		outputs:  map[string]*tasks.CommandCapturer{},
	}

	model.setupDefaultStatus()
	model.setupInterruptHandler()

	p := tea.NewProgram(&model)
	if _, err := p.Run(); err != nil {
		return err
	}

	// TODO: done -display cancel or error or success
	return nil
}

func (m *TreeProgressView) Init() tea.Cmd {
	m.provider.Start()
	return tea.Batch(m.listenToUpdates)
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
	// case spinner.TickMsg:
	// 	var cmd tea.Cmd
	// 	for i := 0; i < len(m.tasksState); i++ {
	// 		task := &m.tasksState[i]
	// 		if task.spinner.ID() == msg.ID {
	// 			task.spinner, cmd = task.spinner.Update(msg)
	// 		}
	// 	}
	// 	return m, cmd
	case TreeProgressUpdateMsg:
		m.mu.Lock()
		m.status[msg.NodeID] = msg.Status
		m.mu.Unlock()
	}

	return m, nil
}

func (m *TreeProgressView) listenToUpdates() tea.Msg {
	return TreeProgressUpdateMsg(<-m.provider.status)
}

func (m *TreeProgressView) View() (s string) {
	s += zwoocBranding + "\n"
	for _, tree := range m.tasks {
		s += tree.Name + "\n"
		tree.Iterate(func(node *tasks.TaskTreeNode) {
			s += fmt.Sprintf("  %s -> %d \n", node.Name, m.status[node.NodeID()])
		})
	}
	return
}

func (m *TreeProgressView) setupDefaultStatus() {
	for _, tree := range m.tasks {
		tree.Iterate(func(node *tasks.TaskTreeNode) {
			// set default status
			m.status[node.NodeID()] = StatusPending
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
