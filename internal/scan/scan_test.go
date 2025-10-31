package scan

import (
	"io/fs"
	"sync"
	"testing"
)

type finalMap struct {
	ExtensionsCount map[string]int
	Mutex           sync.Mutex
}

type fakeDirEntry struct {
	Names         []string
	currentNumber int
}

func (i *fakeDirEntry) Name() string {

	defer func() {
		i.currentNumber++
	}()

	return i.Names[i.currentNumber]
}

func (i *fakeDirEntry) IsDir() bool {
	return false
}

func (i *fakeDirEntry) Type() uint32 {
	return 0
}

func (i *fakeDirEntry) Info() (fs.FileInfo, error) {
	return nil, nil
}

type TestOSInterface struct{}

func (i TestOSInterface) ReadDir(name string) ([]fs.DirEntry, error) {

	const filesCount = 10
	var dirEntry []fakeDirEntry

	return dirEntry, nil
}

func TestScan(t *testing.T) {

	const maxGoroutines = 1000
	finalMap := finalMap{ExtensionsCount: make(map[string]int), Mutex: sync.Mutex{}}

	extensionScanner := NewExtensionScanner(maxGoroutines, TestOSInterface{})
	extensionScanner.ScanFolder("./")
	scanMap, err := extensionScanner.GetScanMap()
	if err != nil {
		t.Errorf("Result was incorrect, error has occurred.")
	}

	finalMap.Mutex.Lock()
	defer func() {
		finalMap.Mutex.Unlock()
	}()

	for key, value := range scanMap {
		finalMap.ExtensionsCount[key] += value
	}
}
