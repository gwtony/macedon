package macedon

import (
	"encoding/base64"
	"encoding/json"
)

type Format struct {
	b64		*base64.Encoding
}

func InitFormat() *Format {
	fm := &Format{}
	fm.b64 = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	return fm
}

func (fm *Format) EncodeBase64(raw []byte) ([]byte, error) {
	return []byte(fm.b64.EncodeToString(raw)), nil
}

func (fm *Format) DecodeBase64(raw []byte) ([]byte, error) {
	return fm.b64.DecodeString(string(raw))
}

func (fm *Format) EncodeJson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (fm *Format) DecodeJson(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
