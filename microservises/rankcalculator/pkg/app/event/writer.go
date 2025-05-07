package event

// TODO PUBLISHER
type Writer interface {
	WriteExchange(evt Event) error
}
