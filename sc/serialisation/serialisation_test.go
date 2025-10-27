package serialisation

import (
	"testing"
)

func TestSerializerRoundtrip(t *testing.T) {
	s := NewSerializer()
	s.WriteUint8(0x7)
	s.WriteUint32LE(0xdeadbeef)
	s.WriteUint64LE(0x1122334455667788)
	s.WriteUvarint(0)
	s.WriteUvarint(300)
	s.WriteBytes([]byte{1, 2, 3})
	s.WriteString("hello")
	s.WriteBool(true)
	s.WriteBool(false)

	data := s.Bytes()
	d := NewDeserializer(data)

	if v, err := d.ReadUint8(); err != nil {
		t.Fatalf("ReadUint8 error: %v", err)
	} else if v != 0x7 {
		t.Fatalf("ReadUint8 expected 0x7 got %d", v)
	}

	if v, err := d.ReadUint32LE(); err != nil {
		t.Fatalf("ReadUint32LE error: %v", err)
	} else if v != 0xdeadbeef {
		t.Fatalf("ReadUint32LE mismatch: got 0x%x", v)
	}

	if v, err := d.ReadUint64LE(); err != nil {
		t.Fatalf("ReadUint64LE error: %v", err)
	} else if v != 0x1122334455667788 {
		t.Fatalf("ReadUint64LE mismatch: got 0x%x", v)
	}

	if v, err := d.ReadUvarint(); err != nil {
		t.Fatalf("ReadUvarint(0) error: %v", err)
	} else if v != 0 {
		t.Fatalf("ReadUvarint(0) mismatch: got %d", v)
	}

	if v, err := d.ReadUvarint(); err != nil {
		t.Fatalf("ReadUvarint(300) error: %v", err)
	} else if v != 300 {
		t.Fatalf("ReadUvarint(300) mismatch: got %d", v)
	}

	if b, err := d.ReadBytes(); err != nil {
		t.Fatalf("ReadBytes error: %v", err)
	} else if len(b) != 3 || b[0] != 1 || b[1] != 2 || b[2] != 3 {
		t.Fatalf("ReadBytes mismatch: %v", b)
	}

	if s, err := d.ReadString(); err != nil {
		t.Fatalf("ReadString error: %v", err)
	} else if s != "hello" {
		t.Fatalf("ReadString mismatch: got %q", s)
	}

	if vb, err := d.ReadBool(); err != nil {
		t.Fatalf("ReadBool(true) error: %v", err)
	} else if vb != true {
		t.Fatalf("ReadBool(true) mismatch: got %v", vb)
	}

	if vb, err := d.ReadBool(); err != nil {
		t.Fatalf("ReadBool(false) error: %v", err)
	} else if vb != false {
		t.Fatalf("ReadBool(false) mismatch: got %v", vb)
	}
}

func TestUvarintBoundaries(t *testing.T) {
	values := []uint64{0, 1, 127, 128, 255, 256, 16383, 16384, 1 << 20, (1 << 63) - 1}
	for _, v := range values {
		s := NewSerializer()
		s.WriteUvarint(v)
		d := NewDeserializer(s.Bytes())
		if got, err := d.ReadUvarint(); err != nil {
			t.Fatalf("uvarint roundtrip error for %d: %v", v, err)
		} else if got != v {
			t.Fatalf("uvarint roundtrip mismatch: wrote %d read %d", v, got)
		}
	}
}
