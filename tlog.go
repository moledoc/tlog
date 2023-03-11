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

type safelogs struct {
	logs     []*log
	writesTo io.Writer
	mu       sync.RWMutex
	t        *testing.T
}

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

var (
	tl = make(map[*testing.T]*safelogs)
)

// -- HELPERS

func print(t *testing.T) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	for _, log := range sl.logs {
		fmt.Fprint(sl.writesTo, log)
	}
	sl.logs = []*log{}
}

func register(t *testing.T, wt io.Writer) *safelogs {
	t.Helper()
	tl[t] = &safelogs{writesTo: wt, t: t}
	t.Cleanup(func() {
		if t.Failed() {
			print(t)
		}
	})
	return tl[t]
}

func lnFormat(count int) string {
	s := make([]string, count)
	for i := 0; i < count; i++ {
		s[i] = "%v"
	}
	return strings.Join(s, " ")
}

func WritesTo(t *testing.T, wt io.Writer) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		sl = register(t, wt)
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.writesTo = wt
}

// -- LOGGERS
func Logf(t *testing.T, format string, args ...any) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		sl = register(t, os.Stdout)
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.logs = append(sl.logs, logline(t, format, args...))
}

func Log(t *testing.T, args ...any) {
	t.Helper()
	Logf(t, lnFormat(len(args)), args...)
}

// -- PRINTERS

func Printf(t *testing.T, format string, args ...any) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	PrintfTo(t, sl.writesTo, format, args...)
}

func PrintfTo(t *testing.T, wt io.Writer, format string, args ...any) {
	t.Helper()
	fmt.Fprint(wt, logline(t, format, args...))
}

func Print(t *testing.T, args ...any) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		sl = register(t, os.Stdout)
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	PrintfTo(t, sl.writesTo, lnFormat(len(args)), args...)
}

func PrintTo(t *testing.T, wt io.Writer, format string, args ...any) {
	t.Helper()
	PrintfTo(t, wt, lnFormat(len(args)), args...)
}

// ---------------------
// ---------------------
// sep obj

func New(t *testing.T) *safelogs {
	return register2(t, os.Stdout)

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

func register2(t *testing.T, wt io.Writer) *safelogs {
	t.Helper()
	sl := &safelogs{writesTo: wt, t: t}
	t.Cleanup(func() {
		if t.Failed() {
			sl.print()
		}
	})
	return sl
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

func (sl *safelogs) Printf(format string, args ...any) {
	sl.t.Helper()
	sl.mu.RLock()
	wt := sl.writesTo
	sl.mu.RUnlock()
	sl.PrintfTo(wt, format, args...)
}

func (sl *safelogs) PrintfTo(wt io.Writer, format string, args ...any) {
	sl.t.Helper()
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	fmt.Fprint(wt, logline(sl.t, format, args...).String())
}

func (sl *safelogs) Print(args ...any) {
	sl.t.Helper()
	// sl.mu.Lock()
	// defer sl.mu.Unlock()
	sl.PrintfTo(sl.writesTo, lnFormat(len(args)), args...)
}

func (sl *safelogs) PrintTo(wt io.Writer, format string, args ...any) {
	sl.t.Helper()
	sl.PrintfTo(wt, lnFormat(len(args)), args...)
}
