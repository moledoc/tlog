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

Few examples.

```go
func TestXxx(t *testing.T) {
    tl := tlog.New(t) // outputs logs to os.Stdout
    // ....
}

func TestXxx(t *testing.T) {
    var bs []byte
    buf := bytes.NewBuffer(bs)
    tl := tlog.NewWithWriter(t, buf) // outputs to bytes.Buffer
}

func TestXxx(t *testing.T) {
    f, _ := os.Open("filename")
    tl = tlog.NewWithWriter(t, f) // outptus to opened file
}
```

## TODOs

* add support for `testing.B` (and other types)

## Author

Meelis Utt