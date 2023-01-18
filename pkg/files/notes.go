package files

type FilesRepo struct {
	path string
}

func New(path string) FilesRepo {
	return FilesRepo{path}
}
