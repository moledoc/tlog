package tlog_test

import (
	// "bytes"
	// "fmt"

	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/moledoc/tlog"
)

// TODO: test panic
// TODO: test fail
// TODO: concurrency

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

func getOutputFile() string {
	filename := actualTestResultsFilename
	if record {
		filename = expectedTestResultsFilename
	}
	return filename
}

func setupTestcase(t *testing.T) *tlog.Logger {
	t.Helper()
	filename := getOutputFile()

	f, err := os.OpenFile(filename, os.O_WRONLY, 0750)
	if err != nil {
		t.Fatalf("unable to open file '%v'\n", filename)
	}

	tl := tlog.NewWithWriter(t, f)
	tl.AddCleanupFunc(func() { fmt.Printf("Closed file '%v' with err '%v'\n", f.Name(), f.Close()) })
	return tl
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

// backup example
// func TestXxx(t *testing.T) {
// 	filename := actualTestResultsFilename
// 	if record {
// 		filename = expectedTestResultsFilename
// 	}

// 	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0750)
// 	if err != nil {
// 		f, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0750)
// 		if err != nil {
// 			t.Fatalf("unable to open file '%v'\n", filename)
// 		}
// 	}

// 	tl := tlog.NewWithWriter(t, f)
// 	tl.AddCleanup(func() { fmt.Println("CLOSED file: ", f.Close()) })
// 	tl.Logf("%v", "one")
// 	panic(1)

// 	tl.Logf("%v,%v", "one", "two")
// 	tl.Logf("\t%v,%v,%v", "one", "two", "three")

// 	t.Fail()
// }

// func TestPanic(t *testing.T) {
// 	// t.SkipNow()
// 	// t.Parallel()
// 	os.Open("actual")
// 	tl := tlog.NewWithWriter(t, os.Stdout)
// 	tl.Log("I will panic")
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Fatal("Expected panic, but nobody paniced")
// 		}

// 	}()
// 	panic(1)
// }

// func TestFailed(t *testing.T) {
// 	t.SkipNow()
// 	t.Parallel()
// 	tl := tlog.New(t)
// 	tl.Log("lala")
// 	tl.Println("hahaha")
// 	t.FailNow()
// }

// func TestTryToPanic2(t *testing.T) {
// 	t.SkipNow()
// 	t.Parallel()
// 	s := time.Now()
// 	cnt := 50000

// 	tl := tlog.New(t)

// 	var wg sync.WaitGroup
// 	// var notExpected string
// 	for i := 0; i < cnt; i++ {
// 		l := strconv.Itoa(i) + ","
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			tl.Log(l)
// 		}()
// 	}
// 	wg.Wait()
// 	// fmt.Println(notExpected)
// 	// tlog.Print(t)
// 	tl.Println(time.Since(s))
// }

// func TestTryToPanic3(t *testing.T) {
// 	t.SkipNow()
// 	t.Parallel()
// 	s := time.Now()
// 	cnt := 50000

// 	tl := tlog.New(t)

// 	var wg sync.WaitGroup
// 	// var notExpected string
// 	for i := 0; i < cnt; i++ {
// 		l := strconv.Itoa(i) + ","
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			tl.Log(l)
// 		}()
// 	}
// 	wg.Wait()
// 	// fmt.Println(notExpected)
// 	// tlog.Println(t)
// 	tl.Println(time.Since(s))
// }
