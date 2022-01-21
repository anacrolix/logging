package log

type Handler interface {
	Handle(Msg)
}
