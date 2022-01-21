package log

type NewLogger struct {
	name     string
	handlers []Handler
	parent   *NewLogger
}

func (l *NewLogger) Handle(m Msg) {
	for _, h := range l.handlers {
		h.Handle(m)
	}
	if l.parent != nil {
		l.parent.Handle(m)
	}
}
