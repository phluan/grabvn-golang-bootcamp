package serializer

import (
	"encoding/xml"
)

type XmlSerializer struct {
}

func (xmlSerializer *XmlSerializer) Render(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}
