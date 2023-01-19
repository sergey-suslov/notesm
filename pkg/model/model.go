package model

import (
	"github.com/charmbracelet/bubbles/key"
	listComponent "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sergey-suslov/notesm/pkg/files"
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

	fr files.FilesRepo
}

func New(fr files.FilesRepo) tea.Model {
	files, err := fr.GetFiles()
	if err != nil {
		panic(err)
	}
	items := make([]listComponent.Item, len(files))
	for i, f := range files {
		items[i] = listComponent.Item(Note{f.Name})
	}

	notesList := listComponent.New(items, listComponent.NewDefaultDelegate(), 20, 20)

	notesList.Title = "projects"
	notesList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			Keymap.Create,
			Keymap.Rename,
			Keymap.Delete,
			Keymap.Back,
		}
	}

	return TeaModel{
		mode:       list,
		terminate:  false,
		notesList:  notesList,
		windowSize: tea.WindowSizeMsg{},

		fr: fr,
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
