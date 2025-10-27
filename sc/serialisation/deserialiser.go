package serialisation

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Deserializer reads typed values from a byte slice produced by Serializer.
type Deserializer struct{ r *bytes.Reader }

// NewDeserializer creates a new Deserializer from bytes.
func NewDeserializer(b []byte) *Deserializer { return &Deserializer{r: bytes.NewReader(b)} }

// ReadUint8 reads one byte.
func (d *Deserializer) ReadUint8() (uint8, error) {
	v, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}
	return v, nil
}

// ReadU16 reads a 2-byte little-endian unsigned integer.
func (d *Deserializer) ReadU16() (uint16, error) {
	var tmp [2]byte
	if _, err := io.ReadFull(d.r, tmp[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(tmp[:]), nil
}

// ReadUint32LE reads a 4-byte little-endian unsigned integer.
func (d *Deserializer) ReadUint32LE() (uint32, error) {
	var tmp [4]byte
	if _, err := io.ReadFull(d.r, tmp[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(tmp[:]), nil
}

// ReadUint64LE reads an 8-byte little-endian unsigned integer.
func (d *Deserializer) ReadUint64LE() (uint64, error) {
	var tmp [8]byte
	if _, err := io.ReadFull(d.r, tmp[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(tmp[:]), nil
}

// ReadUvarint reads a base-128 uvarint.
func (d *Deserializer) ReadUvarint() (uint64, error) {
	v, err := binary.ReadUvarint(d.r)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// ReadBytes reads a length-prefixed byte slice.
func (d *Deserializer) ReadBytes() ([]byte, error) {
	l, err := d.ReadUvarint()
	if err != nil {
		return nil, err
	}
	if l == 0 {
		return []byte{}, nil
	}
	out := make([]byte, l)
	if _, err := io.ReadFull(d.r, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ReadBytesLE reads a 4-byte little-endian length prefix and then the bytes.
func (d *Deserializer) ReadBytesLE() ([]byte, error) {
	var tmp [4]byte
	if _, err := io.ReadFull(d.r, tmp[:]); err != nil {
		return nil, err
	}
	l := binary.LittleEndian.Uint32(tmp[:])
	if l == 0 {
		return []byte{}, nil
	}
	out := make([]byte, l)
	if _, err := io.ReadFull(d.r, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ReadRaw reads n bytes without any length prefix. Caller must ensure there are at least n bytes remaining.
func (d *Deserializer) ReadRaw(n int) ([]byte, error) {
	out := make([]byte, n)
	if _, err := io.ReadFull(d.r, out); err != nil {
		return nil, err
	}
	return out, nil
}

// ReadStringLE reads a 4-byte little-endian length-prefixed string.
func (d *Deserializer) ReadStringLE() (string, error) {
	b, err := d.ReadBytesLE()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadString reads a length-prefixed string.
func (d *Deserializer) ReadString() (string, error) {
	b, err := d.ReadBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadBool reads a single byte boolean.
func (d *Deserializer) ReadBool() (bool, error) {
	b, err := d.ReadUint8()
	if err != nil {
		return false, err
	}
	switch b {
	case 0x00:
		return false, nil
	case 0x01:
		return true, nil
	default:
		return false, fmt.Errorf("invalid boolean value: 0x%02x", b)
	}
}
