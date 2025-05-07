package event

type Publisher interface {
	Publish(body []byte) error
	PublishInExchange(evt Event) error
}
