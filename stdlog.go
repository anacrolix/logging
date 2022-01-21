package log

import (
	"io"
	"log"
)

// Deprecated
var (
	Panicf = log.Panicf
	Fatalf = log.Fatalf
	Fatal  = log.Fatal
)

func New(w io.Writer, prefix string, flags int) *Logger {
	return nil
}
