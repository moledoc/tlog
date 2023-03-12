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

// log

type logline struct {
	time     time.Time
	location string
	testName string
	msg      string
}

func log(t *testing.T, format string, args ...any) *logline {
	t.Helper()
	msg := fmt.Sprintf(format, args...)
	var location string
	for i := 0; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok || !strings.Contains(file, "tlog.go") {
			location = fmt.Sprintf("%v:%v", file, line)
			break
		}
	}
	return &logline{
		time:     time.Now(),
		location: location,
		testName: t.Name(),
		msg:      msg,
	}
}

func (l *logline) String() string {
	return fmt.Sprintf(
		"%v %v %v %v\n",
		l.time.UTC().Format("2006-01-02 15:04:05"),
		l.location,
		fmt.Sprintf("[%v]:", l.testName),
		l.msg,
	)
}

// LOGGER
type logger struct {
	logs     []*logline
	writesTo io.Writer
	mu       sync.RWMutex
	t        *testing.T
}

// HELPERS

func lnFormat(count int) string {
	s := make([]string, count)
	for i := 0; i < count; i++ {
		s[i] = "%v"
	}
	return strings.Join(s, " ")
}

func register(t *testing.T, wt io.Writer) *logger {
	t.Helper()
	sl := &logger{writesTo: wt, t: t}
	t.Cleanup(func() {
		if recover() == nil || t.Failed() {
			sl.print()
		}
	})
	return sl
}

func (sl *logger) print() {
	sl.t.Helper()
	sl.mu.Lock()
	defer sl.mu.Unlock()
	for _, log := range sl.logs {
		fmt.Fprint(sl.writesTo, log)
	}
	sl.logs = []*logline{}
}

func (sl *logger) WritesTo(wt io.Writer) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.writesTo = wt
}

// NEW

func NewWithWriter(t *testing.T, wt io.Writer) *logger {
	return register(t, wt)

}

func New(t *testing.T) *logger {
	return NewWithWriter(t, os.Stdout)

}

// LOGGERS

func (sl *logger) Logf(format string, args ...any) {
	sl.t.Helper()
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.logs = append(sl.logs, log(sl.t, format, args...))
}

func (sl *logger) Log(args ...any) {
	sl.t.Helper()
	sl.Logf(lnFormat(len(args)), args...)
}

// PRINTERS

func (sl *logger) Printf(format string, args ...any) (int, error) {
	sl.t.Helper()
	sl.mu.RLock()
	wt := sl.writesTo
	sl.mu.RUnlock()
	return sl.PrintfTo(wt, format, args...)
}

func (sl *logger) PrintfTo(wt io.Writer, format string, args ...any) (int, error) {
	sl.t.Helper()
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return fmt.Fprint(wt, log(sl.t, format, args...).String())
}

func (sl *logger) Print(args ...any) (int, error) {
	sl.t.Helper()
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.PrintfTo(sl.writesTo, lnFormat(len(args)), args...)
}

func (sl *logger) PrintTo(wt io.Writer, format string, args ...any) (int, error) {
	sl.t.Helper()
	return sl.PrintfTo(wt, lnFormat(len(args)), args...)
}
