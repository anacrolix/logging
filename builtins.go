package logging

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"time"
)

type StreamHandler struct {
	W   io.Writer
	Fmt ByteFormatter
}

func (me StreamHandler) Handle(msg Msg) {
	me.W.Write(me.Fmt(msg.Skip(1)))
}

type ByteFormatter func(Msg) []byte

func DefaultFormatter(msg Msg) []byte {
	name := msg.Name
	if name == "" {
		var pc [1]uintptr
		msg.Callers(1, pc[:])
		name = humanPc(pc[0])
	}
	ret := []byte(fmt.Sprintf("- %s %-5s %s\n  %s",
		time.Now().Format("2006-01-02T15:04:05-0700"),
		func() string {
			if level, ok := msg.GetLevel(); ok {
				return level.LogString()
			}
			return "NONE"
		}(),
		name,
		msg.Text(),
	))
	if ret[len(ret)-1] != '\n' {
		ret = append(ret, '\n')
	}
	return ret
}

func humanPc(pc uintptr) string {
	if pc == 0 {
		panic(pc)
	}
	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	// I'm not sure how to extract just the module from this, since the module might contain valid
	// '.'.
	pkg := f.Function
	file := filepath.Base(f.File)
	return fmt.Sprintf("%s %s:%d", pkg, file, f.Line)
}
