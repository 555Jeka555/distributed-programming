package event

type Writer interface {
	WriteExchange(evt Event) error
}
