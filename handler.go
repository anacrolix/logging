package logging

type Handler interface {
	Handle(Msg)
}
