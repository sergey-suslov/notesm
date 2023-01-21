package files

import (
	"io/ioutil"
	"os"
	"path"
	"time"
)

const DEFAULT_FILES_DIR = ".notesm"

type File struct {
	Name    string
	ModTime time.Time
}

type FilesRepo struct {
	path string
}

func New(path string) FilesRepo {
	return FilesRepo{path}
}

func (fr *FilesRepo) SaveNote(name, body string) error {
	return os.WriteFile(path.Join(fr.path, name), []byte(body), 0644)
}

func (fr *FilesRepo) ReadNote(name string) string {
	content, err := ioutil.ReadFile(path.Join(fr.path, name))
	if err != nil {
		panic(err)
	}
	return string(content)
}

func (fr *FilesRepo) DeleteNote(name string) {
	err := os.Remove(path.Join(fr.path, name))
	if err != nil {
		panic(err)
	}
}

func (fr *FilesRepo) CreateDirIfNotExists() error {
	_, err := ioutil.ReadDir(fr.path)
	if err != nil {
		err = os.Mkdir(fr.path, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fr *FilesRepo) GetFiles() ([]File, error) {
	entries, err := ioutil.ReadDir(fr.path)
	if err != nil {
		return nil, err
	}
	files := make([]File, 0)
	for _, v := range entries {
		if !v.IsDir() && v.Name() != "" {
			files = append(files, File{Name: v.Name(), ModTime: v.ModTime()})
		}
	}
	return files, nil
}
