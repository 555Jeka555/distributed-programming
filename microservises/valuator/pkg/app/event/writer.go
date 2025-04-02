package event

type Writer interface {
	Write(body []byte) error
	WriteExchange(evt Event) error
}
