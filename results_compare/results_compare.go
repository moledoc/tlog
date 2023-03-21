// Copyright 2023 Meelis Utt. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"regexp"
)

var (
	testResultsDir              string = "test_results"
	expectedTestResultsFilename string = testResultsDir + "/expected.log"
	actualTestResultsFilename   string = testResultsDir + "/actual.log"
)

func main() {
	// NOTE: separate the go test output from the compare output
	fmt.Printf("\n\n------------------------------------------------------\n\n")

	expectedResultBytes, err := os.ReadFile(expectedTestResultsFilename)
	if err != nil {
		fmt.Printf("[FATAL]: Failed to open file '%v': %v\n", expectedTestResultsFilename, err)
		os.Exit(1)
	}

	actualResultBytes, err := os.ReadFile(actualTestResultsFilename)
	if err != nil {
		fmt.Printf("[FATAL]: Failed to open file '%v': %v\n", actualTestResultsFilename, err)
		os.Exit(1)
	}

	// NOTE: remove timestamps, because those are not comparable
	re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3} .*/[^:0-9{1-3}]")
	expectedResultLines := re.Split(string(expectedResultBytes), -1)
	actualResultLines := re.Split(string(actualResultBytes), -1)

	// NOTE: check if the results are the same length; if not find the missing/extra lines
	expectedResultLinesCount := len(expectedResultLines)
	actualResultLinesCount := len(actualResultLines)

	passed := true
	if expectedResultLinesCount != actualResultLinesCount {
		passed = false
		expectedLinesMap := make(map[string]struct{})
		for _, l := range expectedResultLines {
			expectedLinesMap[l] = struct{}{}
		}
		actualLinesMap := make(map[string]struct{})
		for _, l := range actualResultLines {
			actualLinesMap[l] = struct{}{}

		}
		var missingLines []string
		var extraLines []string
		for _, line := range expectedResultLines {
			if _, ok := actualLinesMap[line]; !ok {
				missingLines = append(missingLines, line)
			}
		}
		for _, line := range actualResultLines {
			if _, ok := expectedLinesMap[line]; !ok {
				extraLines = append(extraLines, line)
			}
		}
		fmt.Printf("Expected '%v' lines in results, but got '%v'\n", expectedResultLinesCount, actualResultLinesCount)
		if len(missingLines) != 0 {
			fmt.Println("These are the missing lines:")
			for _, line := range missingLines {
				fmt.Printf("* %#v\n", line)
			}
		}
		if len(extraLines) != 0 {
			fmt.Println("These are the extra lines:")
			for _, line := range extraLines {
				fmt.Printf("* %#v\n", line)
			}
		}
	} else {
		// NOTE: check that there are no difference between expected and actual lines.
		for i := 0; i < expectedResultLinesCount; i++ {
			expectedResultLinesi := expectedResultLines[i]
			actualResultLinesi := actualResultLines[i]
			if expectedResultLinesi != actualResultLinesi {
				fmt.Printf("Line %v: expected log entry '%v', but got '%v'\n", i, expectedResultLinesi, actualResultLinesi)
				passed = false
			}
		}
	}

	// NOTE: Make result more clearly visible/separated
	fmt.Println("\n------------------------------------------------------")
	if !passed {
		fmt.Println("[FAILURE]: Actual test results didn't match the expected ones")
		os.Exit(1)
	}

	fmt.Println("[SUCCESS]: Actual test results matched the expected ones")
}
