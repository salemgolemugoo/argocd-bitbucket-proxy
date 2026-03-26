package validator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func Validate(body []byte, signature, secret string) error {
	if signature == "" {
		return fmt.Errorf("missing signature")
	}

	if !strings.HasPrefix(signature, "sha256=") {
		return fmt.Errorf("malformed signature: missing sha256= prefix")
	}

	sigHex := strings.TrimPrefix(signature, "sha256=")
	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return fmt.Errorf("malformed signature: %w", err)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := mac.Sum(nil)

	if !hmac.Equal(sigBytes, expected) {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}
