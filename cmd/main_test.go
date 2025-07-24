package main

import (
	"os"
	"sync"
	"testing"
)

const testFolder = "../test/scanFolder"
const falseFolder = ";;/"

func TestScan(t *testing.T) {

	expectedResults := make(map[string]int)
	expectedResults["csv"] = 1
	expectedResults["yaml"] = 1
	expectedResults["acf"] = 1
	expectedResults["obj"] = 1
	expectedResults["cfg"] = 1
	expectedResults["lua"] = 1
	expectedResults["json"] = 1
	expectedResults["ds_store"] = 2
	expectedResults["ini"] = 2
	expectedResults["dat"] = 3
	expectedResults["txt"] = 4
	expectedResults["png"] = 4
	expectedResults["afl"] = 8

	finalMap := finalMap{ExtensionsCount: make(map[string]int), Mutex: sync.Mutex{}}
	err := runScanWorkers([]string{testFolder}, maxGoroutines, &finalMap)
	if err != nil {
		t.Errorf("Result was incorrect, error has occurred from runScanWorkers")
	}

	if len(finalMap.ExtensionsCount) != len(expectedResults) {
		t.Errorf("Result was incorrect, got len: %d, want len: %d.", len(finalMap.ExtensionsCount), len(expectedResults))
	}

	for extension, extensionsCount1 := range finalMap.ExtensionsCount {
		extensionsCount2, exists := expectedResults[extension]
		if !exists || extensionsCount2 != extensionsCount1 {
			t.Errorf("Result was incorrect, not exists or bad extensions count of %s, got: %d, want: %d.", extension, extensionsCount1, extensionsCount2)
		}
	}

}

func TestPanic(t *testing.T) {

	finalMap := finalMap{ExtensionsCount: make(map[string]int), Mutex: sync.Mutex{}}
	err := runScanWorkers([]string{falseFolder}, maxGoroutines, &finalMap)

	if err == nil {
		t.Errorf("Result was incorrect, want error because incorrect folder.")
	}

}

func TestMain(t *testing.T) {
	os.Args = []string{".", testFolder}
	main()
}
