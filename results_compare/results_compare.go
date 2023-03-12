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
	fmt.Println("IDEA: check that expected and actual logs exist.\nRead in the files and drop the first column or ignore it (ie the timestamp). Then compare line by line the files. They should be the same")

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
	re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} ")
	expectedResultLines := re.Split(string(expectedResultBytes), -1)
	actualResultLines := re.Split(string(actualResultBytes), -1)
	expectedResultLinesCount := len(expectedResultLines)
	actualResultLinesCount := len(actualResultLines)
	if expectedResultLinesCount != actualResultLinesCount {
		fmt.Printf("Expected '%v' lines in results, but got '%v'\n", expectedResultLinesCount, actualResultLinesCount)
		os.Exit(1)
	}
	for i := 0; i < expectedResultLinesCount; i++ {
		expectedResultLinesi := expectedResultLines[i]
		actualResultLinesi := actualResultLines[i]
		if expectedResultLinesi != actualResultLinesi {
			fmt.Printf("Line %v: expected log entry '%v', but got '%v'\n", i, expectedResultLinesi, actualResultLinesi)
			os.Exit(1)
		}
	}

	fmt.Println("Actual test results matched the expected ones")
}
