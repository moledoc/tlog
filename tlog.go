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

// LOG

type log struct {
	time     time.Time
	location string
	testName string
	msg      string
}

func logline(t *testing.T, format string, args ...any) *log {
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
	return &log{
		time:     time.Now(),
		location: location,
		testName: t.Name(),
		msg:      msg,
	}
}

func (l *log) String() string {
	return fmt.Sprintf(
		"%v %v %v %v\n",
		l.time.UTC().Format("2006-01-02 15:04:05"),
		l.location,
		fmt.Sprintf("[%v]:", l.testName),
		l.msg,
	)
}

// LOGGER
type safelogs struct {
	logs     []*log
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

func register(t *testing.T, wt io.Writer) *safelogs {
	t.Helper()
	sl := &safelogs{writesTo: wt, t: t}
	t.Cleanup(func() {
		if t.Failed() {
			sl.print()
		}
	})
	return sl
}

func (sl *safelogs) print() {
	sl.t.Helper()
	sl.mu.Lock()
	defer sl.mu.Unlock()
	for _, log := range sl.logs {
		fmt.Fprint(sl.writesTo, log)
	}
	sl.logs = []*log{}
}

func (sl *safelogs) WritesTo(wt io.Writer) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.writesTo = wt
}

// NEW

func NewWithWriter(t *testing.T, wt io.Writer) *safelogs {
	return register(t, wt)

}

func New(t *testing.T) *safelogs {
	return NewWithWriter(t, os.Stdout)

}

// LOGGERS

func (sl *safelogs) Logf(format string, args ...any) {
	sl.t.Helper()
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.logs = append(sl.logs, logline(sl.t, format, args...))
}

func (sl *safelogs) Log(args ...any) {
	sl.t.Helper()
	sl.Logf(lnFormat(len(args)), args...)
}

// PRINTERS

func (sl *safelogs) Printf(format string, args ...any) (int, error) {
	sl.t.Helper()
	sl.mu.RLock()
	wt := sl.writesTo
	sl.mu.RUnlock()
	return sl.PrintfTo(wt, format, args...)
}

func (sl *safelogs) PrintfTo(wt io.Writer, format string, args ...any) (int, error) {
	sl.t.Helper()
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return fmt.Fprint(wt, logline(sl.t, format, args...).String())
}

func (sl *safelogs) Print(args ...any) (int, error) {
	sl.t.Helper()
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.PrintfTo(sl.writesTo, lnFormat(len(args)), args...)
}

func (sl *safelogs) PrintTo(wt io.Writer, format string, args ...any) (int, error) {
	sl.t.Helper()
	return sl.PrintfTo(wt, lnFormat(len(args)), args...)
}
