package serialisation

import (
	"testing"
)

// helper to assert panic
func expectPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	f()
}

func TestMustNextPanicOnEmpty(t *testing.T) {
	a := NewArgs(nil)

	// Any MustNext should panic because there's no data
	expectPanic(t, func() { _ = a.MustNextU32() })
	expectPanic(t, func() { _ = a.MustNextString() })
	expectPanic(t, func() { _ = a.MustNextU128() })
}

func TestMustNextSuccess(t *testing.T) {
	// Construct args with known data: U32 = 0x0A 00 00 00, string "x"
	a := NewArgs(nil)
	a.AddU32(10)
	a.AddString("x")

	// Read back using the Must helpers (should not panic)
	v := a.MustNextU32()
	if v != 10 {
		t.Fatalf("expected 10, got %d", v)
	}
	s := a.MustNextString()
	if s != "x" {
		t.Fatalf("expected 'x', got '%s'", s)
	}
}
