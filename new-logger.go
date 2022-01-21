package log

import (
	"sync"

	"github.com/anacrolix/torrent/generics"
	"github.com/anacrolix/torrent/option"
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
		h.Handle(m.Skip(1))
	}
	if l.Propagate {
		l.parent.Handle(m.Skip(1))
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

func (l *NewLogger) IsEnabledFor(level Level) bool {
	if l.FilterLevel != NotSet {
		return !level.LessThan(l.FilterLevel)
	}
	if l.parent != nil {
		return l.parent.IsEnabledFor(level)
	}
	return true
}

func (l *NewLogger) LazyLog(level Level, f func() Msg) {
	if l.IsEnabledFor(level) {
		l.Handle(f())
	}
}

func (l *NewLogger) LogLevel(level Level) (ret option.T[ResolvedLogger]) {
	if l.IsEnabledFor(level) {
		return option.Some(ResolvedLogger{})
	}
	return
}

type ResolvedLogger struct {
	l     *NewLogger
	level Level
}

func (me ResolvedLogger) Log(m Msg) {
	m.Level = me.level
	me.l.Handle(m)
}

func (me ResolvedLogger) Logf(format string, args ...interface{}) {
	me.l.Handle(Fstr(format, args...).SetLevel(me.level).Skip(1))
}
