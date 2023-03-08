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
	mu       sync.Mutex
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
	for i := 0; i < 10; i++ {
		_, file, line, _ := runtime.Caller(i)
		if !strings.Contains(file, "tlog.go") {
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
	// MAYBE: use this for time: t.Format(time.UnixDate)
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

func Logf(t *testing.T, format string, args ...any) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		tl[t] = &safelogs{writesTo: os.Stdout}
		sl = tl[t]
	}
	sl.logs = append(sl.logs, logline(t, format, args...))
}

func Print(t *testing.T) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	for _, log := range sl.logs {
		fmt.Fprint(sl.writesTo, log)
	}
}

func LogfPrint(t *testing.T, format string, args ...any) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	Logf(t, format, args...)
	l := sl.logs
	p := l[len(l)-1]
	// TODO: check if p is correct
	fmt.Fprint(sl.writesTo, p)
}

func PrintTo(t *testing.T, wt io.Writer) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	defer Print(t)
	sl.writesTo = wt
}

func lnFormat(count int) string {
	s := make([]string, count)
	for i := 0; i < count; i++ {
		s[i] = "%v"
	}
	return strings.Join(s, " ")
}

func Log(t *testing.T, args ...any) {
	t.Helper()
	Logf(t, lnFormat(len(args)), args...)
}

func LogPrint(t *testing.T, args ...any) {
	t.Helper()
	LogfPrint(t, lnFormat(len(args)), args...)
}

// ----------------

func SafeLogf(t *testing.T, format string, args ...any) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	Logf(t, format, args...)
}

func SafePrint(t *testing.T) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	Print(t)
}

func SafeLogfPrint(t *testing.T, format string, args ...any) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	LogfPrint(t, format, args...)
}

func SafeLog(t *testing.T, args ...any) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	Log(t, args...)
}

func SafeLogPrint(t *testing.T, args ...any) {
	t.Helper()
	sl, ok := tl[t]
	if !ok {
		return
	}
	sl.mu.Lock()
	defer sl.mu.Unlock()
	LogPrint(t, args...)
}
