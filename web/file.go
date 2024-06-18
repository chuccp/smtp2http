package web

import "path/filepath"

type File struct {
	path     string
	fileName string
}

func (f *File) GetPath() string {
	return f.path
}
func (f *File) GetFilename() string {
	if len(f.fileName) > 0 {
		return f.fileName
	}
	return filepath.Base(f.path)
}
