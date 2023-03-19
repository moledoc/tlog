# tlog

Package `tlog` provides logging tools to go along with the normal Go testing system.

The core concept of this package is to enable getting logs for the failed or panicked tests, while not outputting the logs from passing tests to keep the testing logs/output clean and relevant.

It implements types 
* `Entry`, that contains logged info and its metadata
* `Logger` with methods to store and output log entries.

The main `Logger` methods are Log(f) and Print\[f|ln\](To).
The package `tlog` is centered around Log(f) method: this creates a log entry and outputs it only when test fails or panics.
The Print\[f|ln\](To) methods are provided to enable printing the log entry right away without storing it to be printed later.

The logging format in either Log(f) and Print\[f|ln\](To) are uniform and non-configurable, although having some caveats (see Log, Println and PrintlnTo method documentation).

In addition to mentioned, some extra methods are defined to

* define functions that should run after logs are outputted;
* get existing log entries to do additional log parsing manual inside the test;
* mark test as 'panicked', if test itself recovers from the panic;
* change `io.Writer` implementation, to be able to change where the logs are written during the test.

## Usage

In each test it's expected to create a new `Logger` object, using the `New(*testing.T)` or `NewWithWriter(*testing.T, io.Writer)` function.
That logger object can then be used to make log entries to be shown when the test fails/panics (or other actions mentioned above).

Few examples.

```go
func TestXxx(t *testing.T) {
    tl := tlog.New(t) // outputs logs to os.Stdout
	tl.Log("Hello world")
    // ....
}

func TestXxx(t *testing.T) {
    f, _ := os.Open("filename")
    tl = tlog.NewWithWriter(t, f) // outptus to opened file
    tl.AddCleanupFunc(func() { f.Close() }) // close opened file
	tl.Logf("\t%v\n","Hello world")
    // ... 
}


func TestXxx(t *testing.T) {
	tl := tlog.New(t) // outputs to os.Stdout
	tl.Log("Hello")
	time.Sleep(100 * time.Millisecond)
	tl.Log("world")
	entries := tl.GetLogEntries() // get log entries for processing: calculate time between entries
	var timeDiffs []time.Duration
	for i := 1; i < len(entries); i++ {
		timeDiffs = append(timeDiffs, entries[i].Time.Sub(entries[i-1].Time))
	}
	tl.Println("Time differences between log calls:", timeDiffs)
	// ...
}

func TestXxx(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	tl := tlog.NewWithWriter(t, buf) // outputs to bytes.Buffer
	tl.AddCleanupFunc(func() {       // post processing the entries from buffer: calculate time between entries
		fmt.Println(buf.String())
		var timestamps []time.Time
		re := regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{3}")
		lines := strings.Split(buf.String(), "\n")
		for _, line := range lines {
			tsStr := re.FindStringSubmatch(line)
			if len(tsStr) == 0 {
				continue
			}
			ts, _ := time.Parse("2006-01-02 15:04:05", tsStr[0])
			timestamps = append(timestamps, ts)
		}
		var timeDiffs []time.Duration
		for i := 1; i < len(timestamps); i++ {
			timeDiffs = append(timeDiffs, timestamps[i].Sub(timestamps[i-1]))
		}
		fmt.Println("Time differences between log calls:", timeDiffs)
	})
	tl.Log("Hello")
	time.Sleep(100 * time.Millisecond)
	tl.Log("world")
	// ...
	t.FailNow()
}
```

For other examples, see `tlog_test.go` file.

## TODOs

* add support for `testing.B` (and other types)

## Author

Meelis Utt