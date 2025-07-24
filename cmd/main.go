package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/SashaVolohov/scanExtensions/internal/scan"
	"golang.org/x/sync/errgroup"
)

type finalMap struct {
	ExtensionsCount map[string]int
	Mutex           sync.Mutex
}

const noArguments = 1
const maxGoroutines = 1000

func main() {

	if len(os.Args) <= noArguments {
		fmt.Println("This program displays the number of files with a certain extension in the specified folders.")
		fmt.Println("Usage: scanext [folder]...")
		return
	}

	start := time.Now()

	finalMap := finalMap{ExtensionsCount: make(map[string]int), Mutex: sync.Mutex{}}
	runScanWorkers(os.Args[1:], maxGoroutines, &finalMap)

	extensions := slices.Sorted(maps.Keys(finalMap.ExtensionsCount))
	for _, extension := range extensions {
		fmt.Printf("%s: %d\n", extension, finalMap.ExtensionsCount[extension])
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

}

func runScanWorkers(paths []string, maxGoroutines uint, finalMap *finalMap) (err error) {

	group := new(errgroup.Group)

	for i := range paths {

		group.Go(func() error {
			extensionScanner := scan.NewExtensionScanner(maxGoroutines)

			extensionScanner.ScanFolder(paths[i])
			scanMap, err := extensionScanner.GetScanMap()
			if err != nil {
				return err
			}

			finalMap.Mutex.Lock()
			defer func() {
				finalMap.Mutex.Unlock()
			}()

			for key, value := range scanMap {
				finalMap.ExtensionsCount[key] += value
			}

			return nil

		})

	}

	err = group.Wait()
	return
}
