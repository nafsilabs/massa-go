package client

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestReadOnlyResultToBytes_TopLevelNumericArray(t *testing.T) {
	r := &ReadOnlyResult{}
	r.Result.Ok = []interface{}{14.0, 0.0, 0.0, 0.0}

	b, err := ReadOnlyResultToBytes(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(b, []byte{14, 0, 0, 0}) {
		t.Fatalf("unexpected bytes: %v", b)
	}
}

func TestReadOnlyResultToBytes_NestedNumericArray(t *testing.T) {
	r := &ReadOnlyResult{}
	r.Result.Ok = []interface{}{[]interface{}{14.0, 0.0, 0.0, 0.0}}

	b, err := ReadOnlyResultToBytes(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(b, []byte{14, 0, 0, 0}) {
		t.Fatalf("unexpected bytes: %v", b)
	}
}

func TestReadOnlyResultToBytes_Base64String(t *testing.T) {
	payload := []byte{1, 2, 3, 4}
	enc := base64.StdEncoding.EncodeToString(payload)
	r := &ReadOnlyResult{}
	r.Result.Ok = []interface{}{enc}

	b, err := ReadOnlyResultToBytes(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(b, payload) {
		t.Fatalf("unexpected bytes: %v", b)
	}
}

func TestReadOnlyResultToBytes_RawString(t *testing.T) {
	r := &ReadOnlyResult{}
	r.Result.Ok = []interface{}{"hello"}

	b, err := ReadOnlyResultToBytes(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(b, []byte("hello")) {
		t.Fatalf("unexpected bytes: %v", b)
	}
}

func TestReadOnlyResultToBytes_ByteSlice(t *testing.T) {
	r := &ReadOnlyResult{}
	r.Result.Ok = []interface{}{[]byte{9, 8, 7}}

	b, err := ReadOnlyResultToBytes(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(b, []byte{9, 8, 7}) {
		t.Fatalf("unexpected bytes: %v", b)
	}
}

func TestReadOnlyResultToBytes_Unsupported(t *testing.T) {
	r := &ReadOnlyResult{}
	r.Result.Ok = []interface{}{map[string]interface{}{"x": 1}}

	_, err := ReadOnlyResultToBytes(r)
	if err == nil {
		t.Fatalf("expected error for unsupported shape")
	}
}
