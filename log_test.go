package logging

import (
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogBadString(t *testing.T) {
	Str("\xef\x00\xaa\x1ctest\x00test").AddValue("\x00").Log(Default)
}

type stringer struct {
	s string
}

func (me stringer) String() string {
	return strconv.Quote(me.s)
}

func TestValueStringNonLatin(t *testing.T) {
	const (
		u = "カワキヲアメク\n"
		q = `"カワキヲアメク\n"`
	)
	s := stringer{u}
	assert.Equal(t, q, s.String())
	m := Str("").AddValue(q)
	assert.True(t, m.HasValue(q))
}

func BenchmarkDiscardPrintf(b *testing.B) {
	l := GetLogger(b.Name())
	l.Propagate = false
	l.SetHandler(Discard)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Printf("%p: %v", b, i)
	}
}

func BenchmarkFilteredLog(b *testing.B) {
	l := GetLogger(b.Name())
	l.Propagate = false
	l.SetHandler(StreamHandler{
		W:   io.Discard,
		Fmt: DefaultFormatter,
	})
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Printf("%p: %v", b, i)
	}
}
