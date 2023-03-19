// Copyright 2023 Meelis Utt. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// nestedpkg shows that tlog provides proper relative paths for nested packages.
// The results are stored in the same files as tlog_test to be able to run comparison program on them.
package nestedpkg

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/moledoc/tlog"
)

func add(a int, b int) int {
	return a + b
}

var (
	record                      bool
	testResultsDir              string = "../test_results"
	expectedTestResultsFilename string = testResultsDir + "/expected.log"
	actualTestResultsFilename   string = testResultsDir + "/actual.log"
)

func TestMain(m *testing.M) {
	time.Sleep(10 * time.Second) // FIXME: quick and naive fix to make sure nested pkg tests ran after the main ones.
	// NOTE: Parse new flags
	flag.BoolVar(&record, "record", false, "Indicates whether to record new test results or not")
	flag.Parse()

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

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0750)
	if err != nil {
		t.Fatalf("unable to open file '%v'\n", filename)
	}

	tl := tlog.NewWithWriter(t, f)
	tl.AddCleanupFunc(func() { fmt.Printf("Closed file '%v' with err '%v'\n", f.Name(), f.Close()) })
	return tl, f
}

func TestNestedFilepath(t *testing.T) {
	tl, _ := setupTestcase(t)
	a := 34
	b := 35
	tl.Logf("Add(%v,%v)=%v", a, b, add(a, b))
	t.FailNow()
}
