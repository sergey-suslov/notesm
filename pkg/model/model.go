package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	listComponent "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/sergey-suslov/notesm/pkg/files"
	"github.com/sergey-suslov/notesm/pkg/util"
)

type mode int

const (
	list mode = iota
	edit
	detail
	createNoteName
	createNoteBody
	editNoteBody
)

type TeaModel struct {
	notesList  listComponent.Model
	mode       mode
	terminate  bool
	windowSize tea.WindowSizeMsg

	newNoteNameInut textinput.Model
	newNoteName     string

	detail         viewport.Model
	detailNoteName string
	detailHelp     help.Model

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

	newNoteNameInut := textinput.New()
	newNoteNameInut.Placeholder = "Note name"
	newNoteNameInut.CharLimit = 120
	newNoteNameInut.Width = 20

	detail := viewport.New(10, 10)
	detail.YPosition = 0
	detail.MouseWheelEnabled = true

	return TeaModel{
		mode:       list,
		terminate:  false,
		notesList:  notesList,
		windowSize: tea.WindowSizeMsg{},

		newNoteNameInut: newNoteNameInut,

		detail: detail,

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

		headerHeight := lipgloss.Height(TitleStyle.Render(""))
		m.detail.YPosition = headerHeight + 1
		m.detail.Width = msg.Width - v
		m.detail.Height = msg.Height - v*2 - headerHeight - 1

	case abortNoteCreationMsg:
		m.mode = list
		m.newNoteNameInut.Blur()
		m.newNoteNameInut.Reset()
	case newNoteNameResultMsg:
		m.newNoteName = string(msg)
		m.newNoteNameInut.Blur()
		m.newNoteNameInut.Reset()
		cmds = append(cmds, newNoteOpenEditorCmd())
	case updateNotesListMsg:
		cmds = append(cmds, m.notesList.SetItems(m.getNotesAsItems()))
	case editNoteBodyResultMsg:
		m.mode = list
		m.fr.SaveNote(msg.name, msg.content)
	case newNoteBodyResultMsg:
		err := m.createNote(m.newNoteName, string(msg))
		if err != nil {
			panic(err)
		}
		m.mode = list
		cmds = append(cmds, updateNotesListCmd())
	case tea.KeyMsg:
		switch m.mode {
		case detail:
			switch {
			case key.Matches(msg, Keymap.Back):
				m.mode = list
			default:
				m.detail, cmd = m.detail.Update(msg)
				cmds = append(cmds, cmd)
			}
		case createNoteName:
			switch {
			case key.Matches(msg, Keymap.Back):
				cmds = append(cmds, abortNoteCreationCmd())
			case key.Matches(msg, Keymap.Enter):
				cmds = append(cmds, noteNameCreatedCmd(m.newNoteNameInut.Value()))
			default:
				m.newNoteNameInut, cmd = m.newNoteNameInut.Update(msg)
			}
		case list:
			switch {
			case key.Matches(msg, Keymap.Enter):
				m.setDetailContent()
				m.mode = detail
			case key.Matches(msg, Keymap.Delete):
				selected := m.notesList.SelectedItem()
				m.fr.DeleteNote(selected.FilterValue())
				cmds = append(cmds, updateNotesListCmd())
			case key.Matches(msg, Keymap.Edit):
				m.mode = editNoteBody
				selected := m.notesList.SelectedItem()
				noteName := selected.FilterValue()
				cmds = append(cmds, editNoteOpenEditorCmd(noteName, m.fr.ReadNote(noteName)))
			case key.Matches(msg, Keymap.Create):
				m.mode = createNoteName
				cmds = append(cmds, m.newNoteNameInut.Focus(), textinput.Blink)
			default:
				m.notesList, cmd = m.notesList.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
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
	if m.mode == createNoteName {
		return BodyStyle.Render(m.newNoteNameInut.View())
	}
	if m.mode == detail {
		formatted := fmt.Sprintf("%s\n%s", m.noteHeaderView(), m.noteContentView())
		return BodyStyle.Render(formatted)
	}
	return ""
}

func (m TeaModel) defaultBodyMargin() (h, v int) {
	h, v = BodyStyle.GetFrameSize()
	return h, v
}

func (m TeaModel) noteHeaderView() string {
	title := TitleStyle.Render(m.detailNoteName)
	return lipgloss.JoinHorizontal(lipgloss.Center, title)
}

func (m *TeaModel) setDetailContent() {
	selected := m.notesList.SelectedItem()
	noteName := selected.FilterValue()
	m.detailNoteName = noteName
	content := m.fr.ReadNote(noteName)
	_, v := m.defaultBodyMargin()

	width := util.Max(0, util.Min(int(180), m.detail.Width-v*2))
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	content, _ = r.Render(content)

	m.detail.SetContent(content)
}

func (m TeaModel) noteContentView() string {
	_, v := m.defaultBodyMargin()

	return NoteContentStyle(m.detail.Width - v*2).Render(m.detail.View())
}

func (m *TeaModel) getNotesAsItems() []listComponent.Item {
	files, err := m.fr.GetFiles()
	if err != nil {
		panic(err)
	}
	items := make([]listComponent.Item, len(files))
	for i, f := range files {
		items[i] = listComponent.Item(Note{f.Name})
	}
	return items
}
