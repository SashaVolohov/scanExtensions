package scanExtensions

import (
	"io/fs"
	"os"
)

type OSInterface struct{}

func (i OSInterface) ReadDir(name string) ([]fs.DirEntry, error) {
	dirEntry, err := os.ReadDir(name)
	return dirEntry, err
}
