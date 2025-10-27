package serialisation

import (
	"math/big"
	"reflect"
	"testing"
)

func TestArgsStringRoundTrip(t *testing.T) {
	a := NewArgs(nil)
	a.AddString("hello world")
	data := a.Serialise()

	b := NewArgs(data)
	s, err := b.NextString()
	if err != nil {
		t.Fatalf("NextString error: %v", err)
	}
	if s != "hello world" {
		t.Fatalf("string mismatch: got %q", s)
	}
}

func TestArgsArrayRoundTrip(t *testing.T) {
	a := NewArgs(nil)
	in := []byte{1, 2, 3, 4, 5}
	a.AddArray(in)
	data := a.Serialise()

	b := NewArgs(data)
	out, err := b.NextArray()
	if err != nil {
		t.Fatalf("NextArray error: %v", err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("array mismatch: got %v want %v", out, in)
	}
}

func TestArgsBigIntRoundTrip(t *testing.T) {
	// test U128
	a := NewArgs(nil)
	v, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10) // 2^128-1
	a.AddU128(v)
	data := a.Serialise()

	b := NewArgs(data)
	got, err := b.NextU128()
	if err != nil {
		t.Fatalf("NextU128 error: %v", err)
	}
	if got.Cmp(v) != 0 {
		t.Fatalf("U128 mismatch: got %v want %v", got, v)
	}

	// test U256 with a smaller value
	c := NewArgs(nil)
	v2, _ := new(big.Int).SetString("12345678901234567890123456789012345678", 10)
	c.AddU256(v2)
	data2 := c.Serialise()

	d := NewArgs(data2)
	got2, err := d.NextU256()
	if err != nil {
		t.Fatalf("NextU256 error: %v", err)
	}
	if got2.Cmp(v2) != 0 {
		t.Fatalf("U256 mismatch: got %v want %v", got2, v2)
	}
}
