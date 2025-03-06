package transport

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

// EncodeToStream marshals the given message into JSON, prefixes it with a little-endian
// uint32 length and writes it to the provided writer.
func EncodeToStream(w io.Writer, message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	length := uint32(len(data))
	if err := binary.Write(w, binary.LittleEndian, length); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// DecodeFromStream reads a little-endian uint32 length-prefixed JSON message from r and unmarshals it into message.
func DecodeFromStream(r io.Reader, message any) error {
	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return err
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return err
	}
	return json.Unmarshal(data, message)
}
