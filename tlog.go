package tlog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// Entry contains fields to construct a log entry.
type Entry struct {
	Time     time.Time // Timestamp when the log entry was made.
	Location string    // Location (<filepath>:<row number>) where the log entry was made. Eg /foo/bar/baz:54.
	Name     string    // Testcase name, ie testing.T.Name().
	Message  string    // Log message.
}

// String returns log entry as a log string.
// The format used is: <timestamp> <location> [<testname>]: <message>
func (l *Entry) String() string {
	return fmt.Sprintf(
		"%v %v %v %v\n",
		l.Time.UTC().Format("2006-01-02 15:04:05"),
		l.Location,
		fmt.Sprintf("[%v]:", l.Name),
		l.Message,
	)
}

// makeEntry is a function that creates new log entry.
func makeEntry(t *testing.T, format string, args ...any) *Entry {
	t.Helper()
	msg := fmt.Sprintf(format, args...)
	var location string
	for i := 0; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		// MAYBE: think about how to handle !ok better
		if !ok || !strings.Contains(file, "tlog.go") {
			location = fmt.Sprintf("%v:%v", file, line)
			break
		}
	}
	return &Entry{
		Time:     time.Now(),
		Location: location,
		Name:     t.Name(),
		Message:  msg,
	}
}

// Logger is an active logging object that stores log entries and outputs them to an io.Writer, when test fails or panics.
// Logger can be used simultaneously from multiple goroutines; it guarantees to serialize log entries to an internal cache.
type Logger struct {
	// filtered and unexported fields
	t            *testing.T
	writesTo     io.Writer
	logs         []*Entry
	mu           sync.RWMutex
	cleanupFuncs []func()
}

// lnFormat creates a format string with `count` number of values.
// Value format used is '%#v' to get the Go-representation of the values, if object is not a primitive.
// Log and Print methods use the corresponding formatted methods.
// In order to call the corresponding formatted methods, we need to provide a format string.
// lnFormat is used to generate that format string.
func lnFormat(count int) string {
	s := make([]string, count)
	for i := 0; i < count; i++ {
		s[i] = "%#v"
	}
	return strings.Join(s, " ")
}

// createLogger makes a new logger and makes sure that log entries are outputted when the test failed or paniced.
func createLogger(t *testing.T, wt io.Writer) *Logger {
	t.Helper()
	sl := &Logger{writesTo: wt, t: t}
	t.Cleanup(func() {
		if recover() == nil || t.Failed() {
			sl.print()
		}
		for _, fn := range sl.cleanupFuncs {
			fn()
		}
	})
	return sl
}

// print outputs the log entries of the logger to io.Writer specified in the logger object.
func (sl *Logger) print() {
	sl.t.Helper()
	sl.mu.Lock()
	defer sl.mu.Unlock()
	for _, log := range sl.logs {
		fmt.Fprint(sl.writesTo, log)
	}
	sl.logs = []*Entry{}
}

// WritesTo sets the loggers io.Writer to the specified one.
func (sl *Logger) WritesTo(wt io.Writer) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.writesTo = wt
}

// NewWithWriter creates a new logger with provided io.Writer.
func NewWithWriter(t *testing.T, wt io.Writer) *Logger {
	return createLogger(t, wt)
}

// New creates a new logger with os.Stdout as the io.Writer.
func New(t *testing.T) *Logger {
	return NewWithWriter(t, os.Stdout)

}

// AddCleanupFunc adds function to list of functions to be run during the cleanup.
func (sl *Logger) AddCleanupFunc(fn func()) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.cleanupFuncs = append(sl.cleanupFuncs, fn)
}

// Logf formats its arguments according to the format, analogous to fmt.Printf, and records the text in a new log entry.
// A final newline is added if not provided.
// The entry is only outputted when the test fails or panics.
func (sl *Logger) Logf(format string, args ...any) {
	sl.t.Helper()
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.logs = append(sl.logs, makeEntry(sl.t, format, args...))
}

// Log formats its arguments in a default format, analogous to fmt.Println and records the text in a new log entry.
// The entry is only outputted when the test fails or panics.
func (sl *Logger) Log(args ...any) {
	sl.t.Helper()
	sl.Logf(lnFormat(len(args)), args...)
}

// Printf formats its arguments according to the format, analogous to Printf, creates a log entry and outputs it to io.Writer specified in the logger.
// It returns the number of bytes written and any write error.
func (sl *Logger) Printf(format string, args ...any) (int, error) {
	sl.t.Helper()
	sl.mu.RLock()
	wt := sl.writesTo
	sl.mu.RUnlock()
	return sl.PrintfTo(wt, format, args...)
}

// Printf formats its arguments according to the format, analogous to Printf, creates a log entry and outputs it to io.Writer specified in the arguments.
// It returns the number of bytes written and any write error.
func (sl *Logger) PrintfTo(wt io.Writer, format string, args ...any) (int, error) {
	sl.t.Helper()
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return fmt.Fprint(wt, makeEntry(sl.t, format, args...).String())
}

// Println formats its arguments according to the format, analogous to Println, creates a log entry and outputs it to io.Writer specified in the logger.
// It returns the number of bytes written and any write error.
func (sl *Logger) Println(args ...any) (int, error) {
	sl.t.Helper()
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.PrintfTo(sl.writesTo, lnFormat(len(args)), args...)
}

// PrintlnTo formats its arguments according to the format, analogous to Println, creates a log entry and outputs it to io.Writer specified in the arguments.
// It returns the number of bytes written and any write error.
func (sl *Logger) PrintlnTo(wt io.Writer, format string, args ...any) (int, error) {
	sl.t.Helper()
	return sl.PrintfTo(wt, lnFormat(len(args)), args...)
}

// GetLogEntries returns list of log entries recorded in the logger.
func (sl *Logger) GetLogEntries() []*Entry {
	return sl.logs
}
