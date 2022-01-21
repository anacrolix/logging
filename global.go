package log

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	DefaultHandler = StreamHandler{
		W:   os.Stderr,
		Fmt: DefaultFormatter,
	}
	root = NewLogger{
		name:     "",
		handlers: []Handler{DefaultHandler},
		parent:   nil,
	}
	Default = Logger{RootLoggerImpl{&root}}
	Discard = StreamHandler{
		W:   ioutil.Discard,
		Fmt: func(Msg) []byte { return nil },
	}
)

// Terminates old Logger.LoggerImpl chain and tranfers handling to a NewLogger.
type RootLoggerImpl struct {
	*NewLogger
}

func (me RootLoggerImpl) Log(m Msg) {
	me.Handle(m.Skip(1))
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

func Println(a ...interface{}) {
	root.Handle(Msg{
		Args:    a,
		Printer: msgPrintln,
		Skip_:   1,
	})
}
