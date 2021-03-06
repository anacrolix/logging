package logging

import (
	"fmt"
	"sync"

	. "github.com/anacrolix/generics"
)

func GetLogger(name string) *Logger {
	return root.GetChild(name)
}

type Logger struct {
	mu           sync.Mutex
	name         string
	handlers     []Handler
	parent       *Logger
	children     map[string]*Logger
	Propagate    bool
	FilterLevel  Level
	DefaultLevel Level
}

func (l *Logger) Handle(m Msg) {
	for _, h := range l.handlers {
		h.Handle(m.Skip(1))
	}
	if l.Propagate {
		l.parent.Handle(m.Skip(1))
	}
}

func (l *Logger) GetChild(name string) *Logger {
	first, rest := splitName(name)
	l.mu.Lock()
	defer l.mu.Unlock()
	child, ok := l.children[first]
	if ok {
		return child
	}
	child = &Logger{
		name:      fmt.Sprintf("%s.%s", l.name, first),
		parent:    l,
		Propagate: true,
	}
	MakeMapIfNilAndSet(&l.children, first, child)
	if rest != "" {
		return child.GetChild(rest)
	}
	return child
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.Logf(l.DefaultLevel, format, args...)
}

func (l *Logger) Logf(level Level, format string, args ...interface{}) {
	l.Handle(Fstr(format, args...).SetLevel(level).withName(l.name))
}

func (l *Logger) SetHandler(h Handler) {
	l.handlers = []Handler{h}
}

func (l *Logger) IsEnabledFor(level Level) bool {
	if l.FilterLevel != NotSet {
		return !level.LessThan(l.FilterLevel)
	}
	if l.parent != nil {
		return l.parent.IsEnabledFor(level)
	}
	return true
}

func (l *Logger) LazyLog(level Level, f func() Msg) {
	if l.IsEnabledFor(level) {
		l.Handle(f())
	}
}

func (l *Logger) LogLevel(level Level) (ret Option[ResolvedLogger]) {
	if l.IsEnabledFor(level) {
		return Some(ResolvedLogger{})
	}
	return
}

func (l *Logger) LevelOk(level Level) (rl ResolvedLogger, ok bool) {
	ok = l.IsEnabledFor(level)
	if !ok {
		return
	}
	rl = ResolvedLogger{
		l:     l,
		level: level,
	}
	return
}

type ResolvedLogger struct {
	l     *Logger
	level Level
}

func (me ResolvedLogger) Log(m Msg) {
	m.Level = me.level
	me.l.Handle(m)
}

func (me ResolvedLogger) Logf(format string, args ...interface{}) {
	me.l.Handle(Fstr(format, args...).SetLevel(me.level).Skip(1))
}

func (l *Logger) Println(a ...interface{}) {
	l.Handle(Msg{
		Args:    a,
		Printer: msgPrintln,
		Skip_:   1,
	}.withName(l.name))
}

func (l *Logger) Print(a ...interface{}) {
	l.Handle(Msg{
		Args:    a,
		Printer: msgPrint,
		Skip_:   1,
	}.withName(l.name))
}

func msgPrintln(m Msg) string {
	s := fmt.Sprintln(m.Args...)
	return s[:len(s)-1]
}

func msgPrint(m Msg) string {
	return fmt.Sprint(m.Args)
}

func (l *Logger) IsZero() bool {
	return l == nil
}
