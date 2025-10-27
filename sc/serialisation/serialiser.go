package serialisation

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/big"
)

// Serializer helps build byte sequences with common Massa-friendly encodings.
type Serializer struct {
	buf *bytes.Buffer
}

// NewSerializer creates a new Serializer.
func NewSerializer() *Serializer {
	return &Serializer{buf: new(bytes.Buffer)}
}

// Bytes returns the serialized bytes.
func (s *Serializer) Bytes() []byte { return s.buf.Bytes() }

// WriteUint8 appends a single byte.
func (s *Serializer) WriteUint8(v uint8) { s.buf.WriteByte(v) }

// WriteUint32LE appends a 4-byte little-endian unsigned integer.
func (s *Serializer) WriteUint32LE(v uint32) {
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], v)
	s.buf.Write(tmp[:])
}

// WriteU16 writes a 2-byte little-endian unsigned integer.
func (s *Serializer) WriteU16(v uint16) {
	var tmp [2]byte
	binary.LittleEndian.PutUint16(tmp[:], v)
	s.buf.Write(tmp[:])
}

// WriteI16 writes a 2-byte little-endian signed integer.
func (s *Serializer) WriteI16(v int16) {
	s.WriteU16(uint16(v))
}

// WriteI32 writes a 4-byte little-endian signed integer.
func (s *Serializer) WriteI32(v int32) {
	s.WriteUint32LE(uint32(v))
}

// WriteU64 writes an 8-byte little-endian unsigned integer.
func (s *Serializer) WriteU64(v uint64) { s.WriteUint64LE(v) }

// WriteI64 writes an 8-byte little-endian signed integer.
func (s *Serializer) WriteI64(v int64) { s.WriteUint64LE(uint64(v)) }

// WriteF32 writes a 4-byte little-endian IEEE-754 float32.
func (s *Serializer) WriteF32(v float32) {
	s.WriteUint32LE(math.Float32bits(v))
}

// WriteF64 writes an 8-byte little-endian IEEE-754 float64.
func (s *Serializer) WriteF64(v float64) {
	s.WriteUint64LE(math.Float64bits(v))
}

// WriteU128 writes a 16-byte little-endian encoding of a non-negative big.Int.
func (s *Serializer) WriteU128(v *big.Int) {
	b := bigIntToFixedLen(v, 16)
	s.WriteRaw(b)
}

// WriteU256 writes a 32-byte little-endian encoding of a non-negative big.Int.
func (s *Serializer) WriteU256(v *big.Int) {
	b := bigIntToFixedLen(v, 32)
	s.WriteRaw(b)
}

// WriteArrayBytes writes a 4-byte little-endian length-prefixed byte array (used by Dart Args.addArray).
func (s *Serializer) WriteArrayBytes(b []byte) {
	s.WriteBytesLE(b)
}

// bigIntToFixedLen encodes a non-negative big.Int into a fixed-length little-endian byte slice.
func bigIntToFixedLen(v *big.Int, outLen int) []byte {
	if v == nil {
		v = big.NewInt(0)
	}
	if v.Sign() < 0 {
		panic("bigIntToFixedLen: negative value")
	}
	// Create a little-endian buffer of outLen bytes
	out := make([]byte, outLen)
	tmp := v.Bytes() // big-endian
	// tmp is big-endian; copy into out in little-endian order
	for i := 0; i < len(tmp) && i < outLen; i++ {
		out[i] = tmp[len(tmp)-1-i]
	}
	return out
}

// WriteUint64LE appends an 8-byte little-endian unsigned integer.
func (s *Serializer) WriteUint64LE(v uint64) {
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], v)
	s.buf.Write(tmp[:])
}

// WriteUvarint writes a uvarint length-prefix (base-128 varint).
func (s *Serializer) WriteUvarint(v uint64) {
	var tmp [10]byte
	n := binary.PutUvarint(tmp[:], v)
	s.buf.Write(tmp[:n])
}

// WriteBytes writes a length-prefixed byte slice (uvarint length + data).
func (s *Serializer) WriteBytes(b []byte) {
	s.WriteUvarint(uint64(len(b)))
	s.buf.Write(b)
}

// WriteBytesLE writes a 4-byte little-endian length prefix followed by the bytes.
// Use this when the protocol expects a uint32 length prefix instead of uvarint.
func (s *Serializer) WriteBytesLE(b []byte) {
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], uint32(len(b)))
	s.buf.Write(tmp[:])
	s.buf.Write(b)
}

// WriteRaw appends raw bytes without any length prefix.
func (s *Serializer) WriteRaw(b []byte) {
	s.buf.Write(b)
}

// WriteString writes a length-prefixed UTF-8 string.
func (s *Serializer) WriteString(str string) { s.WriteBytes([]byte(str)) }

// WriteStringLE writes a string prefixed with a 4-byte little-endian length.
func (s *Serializer) WriteStringLE(str string) { s.WriteBytesLE([]byte(str)) }

// WriteBool writes a boolean as a single byte (0x00 or 0x01).
func (s *Serializer) WriteBool(v bool) {
	if v {
		s.buf.WriteByte(0x01)
	} else {
		s.buf.WriteByte(0x00)
	}
}
