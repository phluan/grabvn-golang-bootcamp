package serializer

import (
	"encoding/json"
)

type JsonSerializer struct{}

func (jsonSerializer *JsonSerializer) Render(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
