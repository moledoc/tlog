package tlog_test

import (
	// "bytes"
	// "fmt"

	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/moledoc/tlog"
)

// func TestLogf(t *testing.T) {
// 	t.Parallel()
// 	tlog.Logf(t, "%v\n", "test")
// 	tlog.Logf(t, "%v", "test1")
// 	tlog.Logf(t, "%v", "test2")
// 	tlog.Logf(t, "%v", "test3")
// 	tlog.Logf(t, "%v", "test4")
// 	tlog.Log(t, "test4", "test2")
// 	tlog.Print(t)
// }

// func TestLogf2(t *testing.T) {
// 	t.Parallel()
// 	tlog.LogfPrint(t, "\t%v\n", "test")
// 	tlog.LogfPrint(t, "%v", "test1")
// 	tlog.Logf(t, "%v", "test2")
// 	tlog.LogfPrint(t, "%v", "test3")
// 	tlog.Logf(t, "%v", "test4")
// 	tlog.LogPrint(t, "test3", "test4")
// 	// tlog.Print(t)
// }

// func TestConcurrency(t *testing.T) {
// 	t.Parallel()
// 	s := time.Now()
// 	cnt := 1000000
// 	var wg sync.WaitGroup
// 	// var notExpected string
// 	for i := 0; i < cnt; i++ {
// 		l := strconv.Itoa(i) + ","
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			tlog.Log(t, l)
// 			// tlog.LogPrint(t,l)
// 		}()
// 		// notExpected += l
// 	}
// 	wg.Wait()
// 	// fmt.Println(notExpected)
// 	// tlog.Print(t)
// 	tlog.LogPrint(t,time.Since(s))
// }

// func TestConcurrency2(t *testing.T) {
// 	t.Parallel()
// 	s := time.Now()
// 	cnt := 1000000
// 	var wg sync.WaitGroup
// 	// var notExpected string
// 	for i := 0; i < cnt; i++ {
// 		l := strconv.Itoa(i) + ","
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			tlog.Log(t, l)
// 		}()
// 		// notExpected += l
// 	}
// 	wg.Wait()
// 	// fmt.Println(notExpected)
// 	// tlog.Print(t)
// 	tlog.LogPrint(t,time.Since(s))
// }

// func TestPerf(t *testing.T) {
// 	t.Parallel()
// 	s := time.Now()
// 	cnt := 100000
// 	var wg sync.WaitGroup
// 	for i := 0; i < cnt; i++ {
// 		l := strconv.Itoa(i) + ","
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			tlog.Log(t, l)
// 		}()
// 	}
// 	wg.Wait()
// 	tlog.LogPrint(t,time.Since(s))

// 	s = time.Now()
// 	for i := 0; i < cnt; i++ {
// 		l := strconv.Itoa(i) + ","
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			tlog.SafeLog(t, l)
// 		}()
// 	}
// 	wg.Wait()
// 	tlog.LogPrint(t,time.Since(s))
// }

func TestPanic(t *testing.T) {
	t.SkipNow()
	t.Parallel()
	tl := tlog.New(t)
	tl.Log("lala")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic, but nobody paniced")
		}
	}()
	panic(1)
}

func TestFailed(t *testing.T) {
	t.SkipNow()
	t.Parallel()
	tl := tlog.New(t)
	tl.Log("lala")
	tl.Println("hahaha")
	t.FailNow()
}

func TestTryToPanic2(t *testing.T) {
	t.SkipNow()
	t.Parallel()
	s := time.Now()
	cnt := 50000

	tl := tlog.New(t)

	var wg sync.WaitGroup
	// var notExpected string
	for i := 0; i < cnt; i++ {
		l := strconv.Itoa(i) + ","
		wg.Add(1)
		go func() {
			defer wg.Done()
			tl.Log(l)
		}()
	}
	wg.Wait()
	// fmt.Println(notExpected)
	// tlog.Print(t)
	tl.Println(time.Since(s))
}

func TestTryToPanic3(t *testing.T) {
	t.SkipNow()
	t.Parallel()
	s := time.Now()
	cnt := 50000

	tl := tlog.New(t)

	var wg sync.WaitGroup
	// var notExpected string
	for i := 0; i < cnt; i++ {
		l := strconv.Itoa(i) + ","
		wg.Add(1)
		go func() {
			defer wg.Done()
			tl.Log(l)
		}()
	}
	wg.Wait()
	// fmt.Println(notExpected)
	// tlog.Println(t)
	tl.Println(time.Since(s))
}
