package model

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

const DEFAULT_EDITOR = "vim"

type newNoteNameResultMsg string

type newNoteBodyResultMsg string

type updateNotesListMsg struct{}

type abortNoteCreationMsg struct{}

func updateNotesListCmd() tea.Cmd {
	return func() tea.Msg { return updateNotesListMsg{} }
}
func abortNoteCreationCmd() tea.Cmd {
	return func() tea.Msg { return abortNoteCreationMsg{} }
}
func noteNameCreatedCmd(name string) tea.Cmd {
	return func() tea.Msg { return newNoteNameResultMsg(name) }
}
func newNoteOpenEditorCmd() tea.Cmd {
	file, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		panic(err)
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DEFAULT_EDITOR
	}
	c := exec.Command(editor, file.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		bytes, err := os.ReadFile(file.Name())
		if err != nil {
			panic(err)
		}

		return newNoteBodyResultMsg(string(bytes))
	})

}
