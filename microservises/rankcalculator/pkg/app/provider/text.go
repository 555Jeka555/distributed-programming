package provider

type TextProvider interface {
	GetTextID(text string) string
}
