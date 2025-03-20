package event

type Writer interface {
	Write(body []byte) ([]byte, error)
}
