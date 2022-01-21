package log

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	root = NewLogger{
		name: "",
		handlers: []Handler{StreamHandler{
			W:   os.Stderr,
			Fmt: LineFormatter,
		}},
		parent: nil,
	}
	Default = Logger{rootLogger{}}
	Discard = StreamHandler{
		W:   ioutil.Discard,
		Fmt: func(Msg) []byte { return nil },
	}
)

type rootLogger struct{}

func (rootLogger) Log(m Msg) {
	root.Handle(m)
}

func Levelf(level Level, format string, a ...interface{}) {
	Default.Log(Fmsg(format, a...).Skip(1).SetLevel(level))
}

func Printf(format string, a ...interface{}) {
	Default.Log(Fmsg(format, a...).Skip(1))
}

// Prints the arguments to the Default Logger.
func Print(a ...interface{}) {
	// TODO: There's no "Print" equivalent constructor for a Msg, and I don't know what I'd call it.
	Str(fmt.Sprint(a...)).Skip(1).Log(Default)
}
