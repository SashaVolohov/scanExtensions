package scan

import "io/fs"

type osInterface interface {
	ReadDir(name string) ([]fs.DirEntry, error)
}
