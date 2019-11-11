package marshaller

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"

	log "github.com/sirupsen/logrus"
)

// Gob ...
type Gob struct {
	Decoder *gob.Decoder
	Encoder *gob.Encoder
}

// NewGob ...
func NewGob() *Gob {
	return &Gob{}
}

// Marshal ...
func (m *Gob) Marshal(data *StateMetadata) (uint64, error) {
	var output uint64
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return 0, err
	}
	err := binary.Read(buf, binary.BigEndian, &output)

	return output, err
}

// Unmarshal ...
func (m *Gob) Unmarshal(data uint64) (StateMetadata, error) {
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, binary.BigEndian, data); err != nil {
		log.Error("Binary conversion error:", err)
		return StateMetadata{}, err
	}
	log.Debug("Buffer bytes:", buf.Bytes())
	var metadata StateMetadata
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&metadata); err != nil {
		log.Error("Decoding error:", err)
		return metadata, err
	}
	return metadata, nil
}
