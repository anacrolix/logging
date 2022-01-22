package logging

import (
	"runtime"

	"github.com/anacrolix/missinggo/iter"
)

func (m Msg) Text() string {
	return m.Printer(m)
}

func (m Msg) Callers(skip int, pc []uintptr) int {
	return runtime.Callers(m.Skip_+skip+2, pc)
}

func (m Msg) Values(cb iter.Callback) {
	for _, v := range m.Values_ {
		if !cb(v) {
			return
		}
	}
}
