package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestSign_format(t *testing.T) {
	sig, ts := Sign("secret", "POST", "/merchant/v1/deposits", []byte(`{"amount":100000}`))
	if ts == "" || len(ts) < 10 {
		t.Fatalf("bad timestamp: %q", ts)
	}
	if !strings.HasPrefix(sig, "sha256=") {
		t.Errorf("missing sha256= prefix: %q", sig)
	}
	if len(strings.TrimPrefix(sig, "sha256=")) != 64 {
		t.Errorf("hex part wrong length: %q", sig)
	}
}

func TestSign_emptyBody(t *testing.T) {
	sig, ts := Sign("secret", "GET", "/merchant/v1/balance", nil)
	if sig == "" || ts == "" {
		t.Fatal("empty result for GET with no body")
	}
}

func TestSign_algorithm(t *testing.T) {
	// Verify signing_str = timestamp\nMETHOD\npath\nbodyHash
	const (
		secret = "test_secret"
		method = "POST"
		path   = "/merchant/v1/deposits"
	)
	body := []byte(`{"amount":100000}`)

	sig, ts := Sign(secret, method, path, body)

	// Reconstruct expected signature using the same timestamp
	bodyHash := fmt.Sprintf("%x", sha256.Sum256(body))
	signingStr := ts + "\n" + method + "\n" + path + "\n" + bodyHash
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingStr))
	want := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if sig != want {
		t.Errorf("signature mismatch\n got:  %s\nwant: %s", sig, want)
	}
}

func TestSign_differentSecrets(t *testing.T) {
	s1, _ := Sign("secret_a", "GET", "/x", nil)
	s2, _ := Sign("secret_b", "GET", "/x", nil)
	if s1 == s2 {
		t.Error("different secrets produced same signature")
	}
}
