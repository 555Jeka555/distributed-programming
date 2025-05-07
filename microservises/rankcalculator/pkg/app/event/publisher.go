package event

// TODO PUBLISHER
type Publisher interface {
	PublishInExchange(evt Event) error
}
