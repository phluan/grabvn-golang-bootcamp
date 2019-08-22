package serializer

type Serializer interface {
	Render(v interface{}) ([]byte, error)
}
