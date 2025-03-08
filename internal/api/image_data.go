package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/png"
)

type ImageData struct {
	Img image.Image
}

func (id ImageData) MarshalJSON() ([]byte, error) {
	if id.Img == nil {
		return json.Marshal(nil)
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, id.Img); err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return json.Marshal(encoded)
}

func (id *ImageData) UnmarshalJSON(b []byte) error {
	var encoded string
	if err := json.Unmarshal(b, &encoded); err != nil {
		return err
	}
	if encoded == "" {
		id.Img = nil
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	id.Img = img
	return nil
}
