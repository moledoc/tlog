package tlog_test

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
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

func setupTestcase(t *testing.T) (*tlog.Logger, *os.File) {
	t.Helper()
	filename := getOutputFilename()

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0750)
	if err != nil {
		t.Fatalf("unable to open file '%v'\n", filename)
	}

	tl := tlog.NewWithWriter(t, f)
	tl.AddCleanupFunc(func() { fmt.Printf("Closed file '%v' with err '%v'\n", f.Name(), f.Close()) })
	return tl, f
}

func setupTestcaseStdout(t *testing.T) *tlog.Logger {
	t.Helper()
	return tlog.New(t)
}

// TestLogsNoFail shouldn't output anything, since test doesn't fail.
func TestLogsNoFail(t *testing.T) {
	tl, f := setupTestcase(t)
	tl.Logf("one")
	tl.Logf("\t%v\n", "one")
	tl.Logf("%#v%T", "one", f)
	tl.Log("one")
	tl.Log("one", "two")
}

// TestLogs should output logged values, since test fails.
func TestLogs(t *testing.T) {
	tl, f := setupTestcase(t)
	tl.Logf("one")
	tl.Logf("\t%v\n", "one")
	tl.Logf("\n%#v%T", "one", f)
	tl.Log("one")
	tl.Log("one", "two")
	t.Fail()
}

// TestPrints should output logged values, regardless if the test fails
func TestPrints(t *testing.T) {
	tl, f := setupTestcase(t)
	tl.Printf("one")
	tl.Printf("two")
	tl.Printf("%v\t\n%v", "one", "two")
	tl.PrintfTo(f, "one")
	tl.PrintfTo(f, "%v\t\n%v", "one", "two")

	tl.Println("one")
	tl.Println("two")
	tl.Println("%v\t\n%v", "one", "two")
	tl.PrintlnTo(f, "one")
	tl.PrintlnTo(f, "one", "two")
	tl.PrintlnTo(f, "%v\t\n%v", "one", "two")
}

// TestPrintsWithFail should output logged values, regardless if the test fails
func TestPrintsWithFail(t *testing.T) {
	tl, f := setupTestcase(t)
	tl.Printf("one")
	tl.Printf("two")
	tl.Printf("%v\t\n%v", "one", "two")
	tl.PrintfTo(f, "one")
	tl.PrintfTo(f, "%v\t\n%v", "one", "two")

	tl.Println("one")
	tl.Println("two")
	tl.Println("%v\t\n%v", "one", "two")
	tl.PrintlnTo(f, "one")
	tl.PrintlnTo(f, "one", "two")
	tl.PrintlnTo(f, "%v\t\n%v", "one", "two")
	t.Fail()
}

// TestPrintsReturnValues should output correct number of written bytes and errors.
func TestPrintsReturnValues(t *testing.T) {
	tl, f := setupTestcase(t)
	var n int
	var err error
	n, err = tl.Printf("one")
	tl.Println(n, err)
	n, err = tl.Printf("two")
	tl.Println(n, err)
	n, err = tl.Printf("%v\t\n%v", "one", "two")
	tl.Println(n, err)
	n, err = tl.PrintfTo(f, "one")
	tl.Println(n, err)
	n, err = tl.PrintfTo(f, "%v\t\n%v", "one", "two")
	tl.Println(n, err)

	n, err = tl.Println("one")
	tl.Println(n, err)
	n, err = tl.Println("two")
	tl.Println(n, err)
	n, err = tl.Println("%v\t\n%v", "one", "two")
	tl.Println(n, err)
	n, err = tl.PrintlnTo(f, "one")
	tl.Println(n, err)
	n, err = tl.PrintlnTo(f, "one", "two")
	tl.Println(n, err)
	n, err = tl.PrintlnTo(f, "%v\t\n%v", "one", "two")
	tl.Println(n, err)
}

// TestPanics should output logged values when the test panics when logger's SetPanic is called.
func TestPanics(t *testing.T) {
	tl, _ := setupTestcase(t)
	tl.Log("panic at testco")
	defer func() {
		tl.SetPanic()
		if r := recover(); r == nil { // NOTE: r==nil is correct, because we expected a panic and if r is nil, then no panic happened
			tl.Println("Expected panic, but nobody paniced")
		}

	}()
	panic(1)
}

// TestPanicsSetPanicNotCalled shouldn't output logged values when the test panics when logger's SetPanic is not called.
func TestPanicsSetPanicNotCalled(t *testing.T) {
	tl, _ := setupTestcase(t)
	tl.Log("panic at testco")
	defer func() {
		if r := recover(); r == nil { // NOTE: r==nil is correct, because we expected a panic and if r is nil, then no panic happened
			tl.Println("Expected panic, but nobody paniced")
		}

	}()
	panic(1)
}

// TestPanicFromSubFunc should output values when the test panics in a sub func and logger's SetPanic is called.
func TestPanicFromSubFunc(t *testing.T) {
	tl, _ := setupTestcase(t)
	tl.Log("panic at sub-testco")
	func() {
		defer func() {
			tl.SetPanic()
			if r := recover(); r == nil { // NOTE: r==nil is correct, because we expected a panic and if r is nil, then no panic happened
				tl.Println("Expected panic, but nobody paniced")
			}

		}()
		var list []int
		_ = list[100]
	}()
}

// TestConcurrencySafety shouldn't output logged values, since there shouldn't be any data races nor invalid concurrenct object accesses.
func TestConcurrencySafety(t *testing.T) {
	// tl, _ := setupTestcase(t)
	tl := setupTestcaseStdout(t)
	cnt := 100000
	var wg sync.WaitGroup
	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tl.Log(i)
		}(i)
	}
	wg.Wait()
}

// TestRaceConditionDuringTest should output logged values, since there's race condition in the test itself, ie range variable is captured by the func literal.
func TestRaceConditionDuringTest(t *testing.T) {
	tl, _ := setupTestcase(t)
	// tl := setupTestcaseStdout(t)
	cnt := 100
	var wg sync.WaitGroup
	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tl.Log(0 / i) // NOTE: divide 0 with i, so we would have deterministic value in the output.
		}()
	}
	wg.Wait()
}
