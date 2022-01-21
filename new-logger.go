package log

import (
	"sync"

	"github.com/anacrolix/torrent/generics"
)

func GetLogger(name string) *NewLogger {
	return root.GetChild(name)
}

type NewLogger struct {
	mu           sync.Mutex
	name         string
	handlers     []Handler
	parent       *NewLogger
	children     map[string]*NewLogger
	Propagate    bool
	FilterLevel  Level
	DefaultLevel Level
}

func (l *NewLogger) Handle(m Msg) {
	for _, h := range l.handlers {
		h.Handle(m)
	}
	if l.Propagate {
		l.parent.Handle(m)
	}
}

func (l *NewLogger) GetChild(name string) *NewLogger {
	first, rest := splitName(name)
	l.mu.Lock()
	defer l.mu.Unlock()
	child, ok := l.children[first]
	if ok {
		return child
	}
	child = &NewLogger{
		name:      first,
		parent:    l,
		Propagate: true,
	}
	generics.MakeMapIfNilAndSet(&l.children, name, child)
	if rest != "" {
		return child.GetChild(rest)
	}
	return child
}

func (l *NewLogger) Printf(format string, args ...interface{}) {
	l.Logf(l.DefaultLevel, format, args...)
}

func (l *NewLogger) Logf(level Level, format string, args ...interface{}) {
	l.Handle(Fstr(format, args...).SetLevel(level))
}

func (l *NewLogger) SetHandler(h Handler) {
	l.handlers = []Handler{h}
}
