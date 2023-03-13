package tlog_test

import (
	// "bytes"
	// "fmt"

	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/moledoc/tlog"
)

// NOTE: tests for this package should be run sequentially, since expected test results should be in deterministic order.

var (
	record                      bool
	testResultsDir              string = "test_results"
	expectedTestResultsFilename string = testResultsDir + "/expected.log"
	actualTestResultsFilename   string = testResultsDir + "/actual.log"
)

func truncateFile(filename string) {
	f, err := os.OpenFile(filename, os.O_TRUNC, 0750)
	if err != nil {
		fmt.Printf("[WARNING]: Failed to truncate file '%v': %v\n", filename, err)
		f, err = os.OpenFile(filename, os.O_CREATE, 0750)
		if err != nil {
			fmt.Printf("[FATAL]: Failed to open file '%v': %v\n", filename, err)
			os.Exit(1)
		}
	}
	f.Close()
}

func TestMain(m *testing.M) {
	// NOTE: Parse new flags
	flag.BoolVar(&record, "record", false, "Indicates whether to record new test results or not")
	flag.Parse()

	// NOTE: check if test_results dir exist
	if _, err := os.Stat(testResultsDir); errors.Is(err, fs.ErrNotExist) {
		if err := os.Mkdir(testResultsDir, 0750); err != nil {
			fmt.Printf("Failed to create '%v' directory: %v\n", testResultsDir, err)
			os.Exit(1)
		}
	}

	// NOTE: empty actual test result file
	truncateFile(actualTestResultsFilename)

	// NOTE: empty expected test result file, if we are recording
	if record {
		truncateFile(expectedTestResultsFilename)
	}

	os.Exit(m.Run())
}

func getOutputFilename() string {
	filename := actualTestResultsFilename
	if record {
		filename = expectedTestResultsFilename
	}
	return filename
}

func setupTestcase(t *testing.T) *tlog.Logger {
	t.Helper()
	filename := getOutputFilename()

	f, err := os.OpenFile(filename, os.O_WRONLY, 0750)
	if err != nil {
		t.Fatalf("unable to open file '%v'\n", filename)
	}

	tl := tlog.NewWithWriter(t, f)
	tl.AddCleanupFunc(func() { fmt.Printf("Closed file '%v' with err '%v'\n", f.Name(), f.Close()) })
	return tl
}

func setupTestcaseStdout(t *testing.T) *tlog.Logger {
	t.Helper()
	return tlog.New(t)
}

func TestNoFail(t *testing.T) {
	setupTestcase(t)
}

func TestFailOne(t *testing.T) {
	tl := setupTestcase(t)
	tl.Logf("%v", "one")
	t.Fail()
}

func TestFailMulti(t *testing.T) {
	tl := setupTestcase(t)
	tl.Logf("%v", "one")
	tl.Logf("%v,%v", "one", "two")
	tl.Logf("\t%v,%v,%v", "one", "two", "three")
	t.Fail()
}

func TestFailMultiMix(t *testing.T) {
	tl := setupTestcase(t)
	tl.Log("one")
	tl.Log("one", "two")
	tl.Logf("\t%v,%v,%v", "one", "two", "three")
	t.Fail()
}

func TestPanic(t *testing.T) {
	tl := setupTestcase(t)
	tl.Log("panic at testco")
	defer func() {
		tl.SetPanic()
		if r := recover(); r == nil {
			t.Fatal("Expected panic, but nobody paniced")
		}

	}()
	panic(1)
}

func TestFailedTestcase(t *testing.T) {
	tl := setupTestcase(t)
	tl.Log("this testcase is a failure")
	t.FailNow()
}

func TestConcurrencySafety(t *testing.T) {
	tl := setupTestcase(t)
	// tl := setupTestcaseStdout(t)
	cnt := 100000
	var wg sync.WaitGroup
	for i := 0; i < cnt; i++ {
		l := strconv.Itoa(i) + ","
		wg.Add(1)
		go func() {
			defer wg.Done()
			tl.Log(l)
		}()
	}
	wg.Wait()
}
