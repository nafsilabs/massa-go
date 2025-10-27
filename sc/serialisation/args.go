package serialisation

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
)

// Args provides a combined serializer/deserializer with an internal offset
// similar to Dart's Args class. It tracks offset and delegates encoding/decoding
// to the Serializer and Deserializer helpers to avoid duplication.
type Args struct {
	ser    *Serializer
	rdr    []byte
	offset int
}

// NewArgs creates an Args builder. If initialData is provided, it will be used as
// the initial buffer for deserialization and offset will be set to 0.
func NewArgs(initialData []byte) *Args {
	a := &Args{ser: NewSerializer(), offset: 0}
	if initialData != nil {
		// initialize serializer buffer with initial data so Serialise() returns it
		a.ser.buf.Write(initialData)
		a.rdr = initialData
	}
	return a
}

// Serialise returns the serialized bytes.
func (a *Args) Serialise() []byte { return a.ser.Bytes() }

// Offset returns the current read offset.
func (a *Args) Offset() int { return a.offset }

// ensureReader makes sure rdr is up to date from serializer buffer when needed.
func (a *Args) ensureReader() {
	if a.rdr == nil || len(a.ser.Bytes()) > len(a.rdr) {
		a.rdr = a.ser.Bytes()
	}
}

// helper to create a Deserializer over the remaining bytes and track consumed bytes
func (a *Args) newDeserializerAtOffset() *Deserializer {
	a.ensureReader()
	// create deserializer over rdr[offset:]
	return &Deserializer{r: bytes.NewReader(a.rdr[a.offset:])}
}

// advance offset by consumed bytes (reader.Len() reports remaining bytes)
func (a *Args) advanceFromReader(d *Deserializer, before int) {
	after := d.r.Len()
	consumed := before - after
	a.offset += consumed
}

// --- Add (serialize) methods ---
func (a *Args) AddBool(v bool)       { a.ser.WriteBool(v) }
func (a *Args) AddU8(v uint8)        { a.ser.WriteUint8(v) }
func (a *Args) AddU16(v uint16)      { a.ser.WriteU16(v) }
func (a *Args) AddI16(v int16)       { a.ser.WriteI16(v) }
func (a *Args) AddU32(v uint32)      { a.ser.WriteUint32LE(v) }
func (a *Args) AddI32(v int32)       { a.ser.WriteI32(v) }
func (a *Args) AddU64(v *big.Int)    { a.ser.WriteU64(v.Uint64()) }
func (a *Args) AddI64(v *big.Int)    { a.ser.WriteI64(int64(v.Int64())) }
func (a *Args) AddU128(v *big.Int)   { a.ser.WriteU128(v) }
func (a *Args) AddU256(v *big.Int)   { a.ser.WriteU256(v) }
func (a *Args) AddF32(v float32)     { a.ser.WriteF32(v) }
func (a *Args) AddF64(v float64)     { a.ser.WriteF64(v) }
func (a *Args) AddString(s string)   { a.ser.WriteStringLE(s) }
func (a *Args) AddArray(data []byte) { a.ser.WriteArrayBytes(data) }

// --- Next (deserialize) methods ---
func (a *Args) NextBool() (bool, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	v, err := d.ReadBool()
	if err != nil {
		return false, fmt.Errorf("NextBool: %w", err)
	}
	a.advanceFromReader(d, before)
	return v, nil
}

func (a *Args) NextU8() (uint8, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	v, err := d.ReadUint8()
	if err != nil {
		return 0, fmt.Errorf("NextU8: %w", err)
	}
	a.advanceFromReader(d, before)
	return v, nil
}

func (a *Args) NextU16() (uint16, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	v, err := d.ReadU16()
	if err != nil {
		return 0, fmt.Errorf("NextU16: %w", err)
	}
	a.advanceFromReader(d, before)
	return v, nil
}

func (a *Args) NextI16() (int16, error) {
	v, err := a.NextU16()
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func (a *Args) NextU32() (uint32, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	v, err := d.ReadUint32LE()
	if err != nil {
		return 0, fmt.Errorf("NextU32: %w", err)
	}
	a.advanceFromReader(d, before)
	return v, nil
}

func (a *Args) NextI32() (int32, error) {
	v, err := a.NextU32()
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func (a *Args) NextU64() (*big.Int, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	v, err := d.ReadUint64LE()
	if err != nil {
		return nil, fmt.Errorf("NextU64: %w", err)
	}
	a.advanceFromReader(d, before)
	return big.NewInt(0).SetUint64(v), nil
}

func (a *Args) NextI64() (*big.Int, error) { return a.NextU64() }

func (a *Args) NextU128() (*big.Int, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	b, err := d.ReadRaw(16)
	if err != nil {
		return nil, fmt.Errorf("NextU128: %w", err)
	}
	a.advanceFromReader(d, before)
	return bytesToBigIntLE(b), nil
}

func (a *Args) NextU256() (*big.Int, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	b, err := d.ReadRaw(32)
	if err != nil {
		return nil, fmt.Errorf("NextU256: %w", err)
	}
	a.advanceFromReader(d, before)
	return bytesToBigIntLE(b), nil
}

func (a *Args) NextF32() (float32, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	v, err := d.ReadUint32LE()
	if err != nil {
		return 0, fmt.Errorf("NextF32: %w", err)
	}
	a.advanceFromReader(d, before)
	return math.Float32frombits(v), nil
}

func (a *Args) NextF64() (float64, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	v, err := d.ReadUint64LE()
	if err != nil {
		return 0, fmt.Errorf("NextF64: %w", err)
	}
	a.advanceFromReader(d, before)
	return math.Float64frombits(v), nil
}

func (a *Args) NextString() (string, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	s, err := d.ReadStringLE()
	if err != nil {
		return "", fmt.Errorf("NextString: %w", err)
	}
	a.advanceFromReader(d, before)
	return s, nil
}

// NextArray returns the raw bytes of the array (length-prefixed by u32)
func (a *Args) NextArray() ([]byte, error) {
	d := a.newDeserializerAtOffset()
	before := d.r.Len()
	b, err := d.ReadBytesLE()
	if err != nil {
		return nil, fmt.Errorf("NextArray: %w", err)
	}
	a.advanceFromReader(d, before)
	return b, nil
}

// helper to convert little-endian bytes to big.Int
func bytesToBigIntLE(b []byte) *big.Int {
	// reverse to big-endian
	be := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		be[i] = b[len(b)-1-i]
	}
	return new(big.Int).SetBytes(be)
}

// --- MustNext helpers: convenience wrappers that panic on error. Useful in
// examples/tests where brevity is preferred over error handling.
func (a *Args) MustNextBool() bool {
	v, err := a.NextBool()
	if err != nil {
		panic(err)
	}
	return v
}

func (a *Args) MustNextU8() uint8 {
	v, err := a.NextU8()
	if err != nil {
		panic(err)
	}
	return v
}

func (a *Args) MustNextU16() uint16 {
	v, err := a.NextU16()
	if err != nil {
		panic(err)
	}
	return v
}

func (a *Args) MustNextU32() uint32 {
	v, err := a.NextU32()
	if err != nil {
		panic(err)
	}
	return v
}

func (a *Args) MustNextU64() *big.Int {
	v, err := a.NextU64()
	if err != nil {
		panic(err)
	}
	return v
}

func (a *Args) MustNextU128() *big.Int {
	v, err := a.NextU128()
	if err != nil {
		panic(err)
	}
	return v
}

func (a *Args) MustNextU256() *big.Int {
	v, err := a.NextU256()
	if err != nil {
		panic(err)
	}
	return v
}

func (a *Args) MustNextString() string {
	v, err := a.NextString()
	if err != nil {
		panic(err)
	}
	return v
}

func (a *Args) MustNextArray() []byte {
	v, err := a.NextArray()
	if err != nil {
		panic(err)
	}
	return v
}
