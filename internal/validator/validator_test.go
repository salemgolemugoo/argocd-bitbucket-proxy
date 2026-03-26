package validator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func computeHMAC(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestValidate_ValidSignature(t *testing.T) {
	body := []byte(`{"test": "payload"}`)
	secret := "my-secret"
	sig := computeHMAC(body, secret)

	err := Validate(body, sig, secret)
	if err != nil {
		t.Fatalf("expected valid signature, got error: %v", err)
	}
}

func TestValidate_InvalidSignature(t *testing.T) {
	body := []byte(`{"test": "payload"}`)
	wrongSig := "sha256=0000000000000000000000000000000000000000000000000000000000000000"

	err := Validate(body, wrongSig, "my-secret")
	if err == nil {
		t.Fatal("expected error for invalid signature")
	}
}

func TestValidate_MissingSignature(t *testing.T) {
	body := []byte(`{"test": "payload"}`)
	err := Validate(body, "", "my-secret")
	if err == nil {
		t.Fatal("expected error for missing signature")
	}
}

func TestValidate_MalformedSignature(t *testing.T) {
	body := []byte(`{"test": "payload"}`)
	err := Validate(body, "not-a-valid-sig", "my-secret")
	if err == nil {
		t.Fatal("expected error for malformed signature")
	}
}
