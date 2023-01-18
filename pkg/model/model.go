package model

import (
	listComponent "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	list mode = iota
	edit
	create
)

type TeaModel struct {
	notesList  listComponent.Model
	mode       mode
	terminate  bool
	windowSize tea.WindowSizeMsg
}

func New() tea.Model {
	items := make([]listComponent.Item, 2)
	for i, proj := range []Note{{"Item 1"}, {"Item 2"}} {
		items[i] = listComponent.Item(proj)
	}

	return TeaModel{
		mode:       list,
		terminate:  false,
		notesList:  listComponent.New(items, listComponent.NewDefaultDelegate(), 20, 20),
		windowSize: tea.WindowSizeMsg{},
	}
}

// Init implements tea.Model
func (TeaModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := BodyStyle.GetFrameSize()
		m.notesList.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		m.notesList, cmd = m.notesList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m TeaModel) View() string {
	if m.terminate {
		return ""
	}
	if m.mode == list {
		return BodyStyle.Render(m.notesList.View())
	}
	return ""
}
