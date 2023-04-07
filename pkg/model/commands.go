package model

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

const DEFAULT_EDITOR = "vim"

type newNoteNameResultMsg string

type newNoteBodyResultMsg string

type editNoteBodyResultMsg struct {
	name    string
	content string
}

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

func editNoteOpenEditorCmd(name string, content string) tea.Cmd {
	file, err := os.CreateTemp(os.TempDir(), "*.md")
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(content)
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
			println(err)
			panic(err)
		}

		return editNoteBodyResultMsg{
			name: name, content: string(bytes),
		}
	})
}

func newNoteOpenEditorCmd() tea.Cmd {
	file, err := os.CreateTemp(os.TempDir(), "*.md")
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
