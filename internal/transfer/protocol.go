package transfer

import (
	"encoding/binary"
	"fmt"
	"io"
)

// WritePacket writes a length-prefixed byte slice over an io.Writer
func WritePacket(w io.Writer, payload []byte) error {
	length := uint32(len(payload))
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return fmt.Errorf("failed to write packet length header: %w", err)
	}

	if _, err := w.Write(payload); err != nil {
		return fmt.Errorf("failed to write packet payload: %w", err)
	}

	return nil
}

// ReadPacket reads a length-prefixed byte slice from an io.Reader
func ReadPacket(r io.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to read packet length header: %w", err)
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, fmt.Errorf("failed to read packet payload body: %w", err)
	}

	return payload, nil
}
