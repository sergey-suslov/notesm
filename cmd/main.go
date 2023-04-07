package main

import (
	"fmt"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sergey-suslov/notesm/pkg/files"
	"github.com/sergey-suslov/notesm/pkg/model"
)

func main() {
	homeDir := os.Getenv("HOME")
	filesRepoPath := path.Join(homeDir, files.DEFAULT_FILES_DIR)
	filesRepo := files.New(filesRepoPath)
	_ = filesRepo.CreateDirIfNotExists()

	p := tea.NewProgram(model.New(filesRepo), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
