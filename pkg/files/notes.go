package files

import (
	"io/ioutil"
	"os"
	"path"
)

const DEFAULT_FILES_DIR = ".notesm"

type File struct {
	Name string
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
			files = append(files, File{v.Name()})
		}
	}
	return files, nil
}
