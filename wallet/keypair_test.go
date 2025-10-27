package wallet

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/ed25519"
)

func TestVersionedSignatureAndParse(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("generate keypair: %v", err)
	}

	message := []byte("hello massa")

	vsig := kp.VersionedSignatureBytes(message)
	if len(vsig) != 1+ed25519.SignatureSize {
		t.Fatalf("unexpected versioned signature length: %d", len(vsig))
	}

	v, sig, err := ParseVersionedSignatureBytes(vsig)
	if err != nil {
		t.Fatalf("parse versioned signature bytes: %v", err)
	}
	if v != 0x00 {
		t.Fatalf("expected version 0x00, got 0x%02x", v)
	}

	// verify raw
	if !kp.Verify(message, sig) {
		t.Fatalf("raw signature verification failed")
	}

	// verify via helper
	if !VerifyVersionedSignature(message, vsig, kp.PublicKey) {
		t.Fatalf("VerifyVersionedSignature failed")
	}

	b58 := kp.VersionedSignatureBase58(message)
	if len(b58) == 0 {
		t.Fatalf("empty base58 signature")
	}

	v2, sig2, err := ParseVersionedSignatureBase58(b58)
	if err != nil {
		t.Fatalf("parse base58: %v", err)
	}
	if v2 != 0x00 {
		t.Fatalf("expected version 0x00, got 0x%02x", v2)
	}
	if !bytes.Equal(sig, sig2) {
		t.Fatalf("parsed signatures differ")
	}

	if !VerifyVersionedSignatureBase58(message, b58, kp.PublicKey) {
		t.Fatalf("VerifyVersionedSignatureBase58 failed")
	}
}
