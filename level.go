package logging

type Level struct {
	rank   int
	logStr string
}

var levelKey = new(struct{})

var (
	NotSet   = Level{0, "NONE"}
	Debug    = Level{1, "DEBUG"}
	Info     = Level{2, "INFO"}
	Warning  = Level{3, "WARN"}
	Error    = Level{4, "ERROR"}
	Critical = Level{5, "CRIT"}
	// Will this get special treatment? Not yet. Also disabled due to conflict with std log.Fatal.
	//Fatal = Level{6, "FATAL"}
)

func (l Level) LogString() string {
	return l.logStr
}

func (l Level) LessThan(r Level) bool {
	return l.rank < r.rank
}
