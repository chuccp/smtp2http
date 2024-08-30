package web

import "path/filepath"

type File struct {
	Path     string
	FileName string
}

func (f *File) GetPath() string {
	return f.Path
}
func (f *File) GetFilename() string {
	if len(f.FileName) > 0 {
		return f.FileName
	}
	return filepath.Base(f.Path)
}
