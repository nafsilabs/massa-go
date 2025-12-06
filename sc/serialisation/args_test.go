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

func TestArgsUint32ArrayComparison(t *testing.T) {
	//first := uint32(10)
	aa := NewArgs(nil)
	//aa.AddU32(first)

	values := []uint32{100, 200, 300, 400, 500}

	// Approach 1: Add individually
	a := NewArgs(nil)
	for _, v := range values {
		a.AddU32(v)
	}
	aa.AddArray(a.Serialise()) //258,000/-
	data1 := aa.Serialise()

	// Approach 2: Use AddU32Array
	b := NewArgs(nil)
	//b.AddU32(first)
	b.AddU32Array(values)
	data2 := b.Serialise()

	// Compare serialized results
	if !reflect.DeepEqual(data1, data2) {
		t.Fatalf("serialized data mismatch:\nindividual: %v\narray:      %v", data1, data2)
	}

	// Verify deserialization from approach 1
	c := NewArgs(data1)
	data1Bytes, err := c.NextArray()
	if err != nil {
		t.Fatalf("NextArray error: %v", err)
	}
	c = NewArgs(data1Bytes)
	for i, want := range values {
		got, err := c.NextU32()
		if err != nil {
			t.Fatalf("NextU32 error at index %d: %v", i, err)
		}
		if got != want {
			t.Fatalf("value mismatch at index %d: got %d want %d", i, got, want)
		}
	}

	// Verify deserialization from approach 2
	d := NewArgs(data2)
	gotArray, err := d.NextU32Array()
	if err != nil {
		t.Fatalf("NextUint32Array error: %v", err)
	}
	if !reflect.DeepEqual(gotArray, values) {
		t.Fatalf("array mismatch: got %v want %v", gotArray, values)
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
