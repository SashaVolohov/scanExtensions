package scan

import (
	"fmt"
	"maps"
	"os"
	"strings"
	"sync"
)

const extensionSeparator = "."
const noActiveWorkers = 0
const withoutFileExtension = 1

type extensionScanner struct {
	extensionsCount map[string]int
	mutex           sync.Mutex
	os              osInterface

	done             chan struct{}
	mainSemaphore    chan struct{}
	workersSemaphore chan struct{}

	err      error
	errMutex sync.Mutex
}

func NewExtensionScanner(maxGoroutines uint, osInterface osInterface) *extensionScanner {
	return &extensionScanner{extensionsCount: make(map[string]int), mutex: sync.Mutex{}, errMutex: sync.Mutex{}, done: make(chan struct{}), workersSemaphore: make(chan struct{}, maxGoroutines), mainSemaphore: make(chan struct{}, 1), os: osInterface}
}

func (e *extensionScanner) checkExtension(path string) {
	splittedName := strings.Split(path, extensionSeparator)
	if len(splittedName) == withoutFileExtension {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.extensionsCount[strings.ToLower(splittedName[len(splittedName)-1])]++
}

func (e *extensionScanner) readDir(path string) []os.DirEntry {
	files, err := e.os.ReadDir(path)
	if err != nil {

		e.errMutex.Lock()
		defer e.errMutex.Unlock()

		if e.err == nil {
			e.err = err
		}
	}
	return files
}

func (e *extensionScanner) scannerWorker(path string) {

	e.workersSemaphore <- struct{}{}

	go func() {

		defer func() {
			<-e.workersSemaphore

			if len(e.workersSemaphore) == noActiveWorkers {
				e.done <- struct{}{}
			}
		}()

		files := e.readDir(path)
		if e.err != nil {
			return
		}

		for _, val := range files {
			if val.IsDir() {
				e.scannerWorker(fmt.Sprintf("%s\\%s", path, val.Name()))
			} else {
				e.checkExtension(val.Name())
			}
		}

	}()

}

func (e *extensionScanner) ScanFolder(path string) {

	e.mainSemaphore <- struct{}{}
	defer func() {
		<-e.mainSemaphore
	}()

	e.scannerWorker(path)
	<-e.done

}

func (e *extensionScanner) GetScanMap() (map[string]int, error) {
	copiedMap := make(map[string]int)
	maps.Copy(copiedMap, e.extensionsCount)

	return copiedMap, e.err
}
