package model

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

var (
	BodyStyle  = lipgloss.NewStyle().Margin(1, 2)
	TitleStyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()

		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1).Margin(0).Bold(true)
	}()
)

var NoteContentStyle = func(width int) lipgloss.Style {
	b := lipgloss.NormalBorder()

	return lipgloss.NewStyle().Width(width).BorderStyle(b).Padding(0, 0)
}

type keymap struct {
	Create    key.Binding
	Enter     key.Binding
	Rename    key.Binding
	Edit      key.Binding
	Delete    key.Binding
	Back      key.Binding
	BackWithQ key.Binding
	Quit      key.Binding
}

var Keymap = keymap{
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	BackWithQ: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}

func (k keymap) NotesListHelp() []key.Binding {
	return []key.Binding{
		Keymap.Create,
		Keymap.Rename,
		Keymap.Delete,
		Keymap.Back,
	}
}
