package surepay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
)

// buildWebhookBody creates a signed webhook body for testing.
func buildWebhookBody(t *testing.T, secret string, fields map[string]any) []byte {
	t.Helper()
	canonical, err := marshalSorted(fields)
	if err != nil {
		t.Fatalf("marshalSorted: %v", err)
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(canonical)
	fields["signature"] = hex.EncodeToString(mac.Sum(nil))
	body, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return body
}

func TestVerifyWebhookSignature_valid(t *testing.T) {
	secret := "webhook_secret_test"
	payload := map[string]any{
		"event":     "deposit.success",
		"amount":    float64(100000),
		"currency":  "VND",
		"status":    "success",
		"timestamp": float64(1704067200),
	}
	body := buildWebhookBody(t, secret, payload)

	if !VerifyWebhookSignature(secret, body, "") {
		t.Error("expected valid signature to pass")
	}
}

func TestVerifyWebhookSignature_wrongSecret(t *testing.T) {
	body := buildWebhookBody(t, "correct", map[string]any{
		"event": "deposit.success", "amount": float64(100000),
	})
	if VerifyWebhookSignature("wrong_secret", body, "") {
		t.Error("wrong secret should fail")
	}
}

func TestVerifyWebhookSignature_tamperedBody(t *testing.T) {
	payload := map[string]any{"event": "payout.success", "amount": float64(500000)}
	body := buildWebhookBody(t, "sec", payload)

	// Tamper: change amount in the body JSON
	var m map[string]any
	json.Unmarshal(body, &m)
	m["amount"] = float64(1) // attacker changes amount
	tampered, _ := json.Marshal(m)

	if VerifyWebhookSignature("sec", tampered, "") {
		t.Error("tampered body should fail")
	}
}

func TestVerifyWebhookSignature_missingSig(t *testing.T) {
	body := []byte(`{"event":"deposit.success","amount":100000}`)
	if VerifyWebhookSignature("sec", body, "") {
		t.Error("missing signature should fail")
	}
}

func TestVerifyWebhookSignature_sigFromParam(t *testing.T) {
	secret := "sec"
	fields := map[string]any{"event": "deposit.success", "amount": float64(100)}
	canonical, _ := marshalSorted(fields)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(canonical)
	sig := hex.EncodeToString(mac.Sum(nil))

	// Body has no "signature" field — sig supplied via parameter (header)
	body, _ := json.Marshal(fields)
	if !VerifyWebhookSignature(secret, body, sig) {
		t.Error("sig from header parameter should pass")
	}
}

func TestVerifyWebhookSignature_keyOrderInvariant(t *testing.T) {
	// Regardless of JSON marshal key order, canonical sorting must produce same result
	secret := "ord_test"
	fields := map[string]any{
		"z_field": "zzz",
		"a_field": "aaa",
		"m_field": float64(123),
	}
	body := buildWebhookBody(t, secret, fields)
	if !VerifyWebhookSignature(secret, body, "") {
		t.Error("key-order invariant failed")
	}
}

func TestVerifyWebhookSignature_invalidJSON(t *testing.T) {
	if VerifyWebhookSignature("sec", []byte("not json"), "") {
		t.Error("invalid JSON should fail")
	}
}

func TestMarshalSorted(t *testing.T) {
	m := map[string]any{"b": "two", "a": "one", "c": float64(3)}
	got, err := marshalSorted(m)
	if err != nil {
		t.Fatalf("marshalSorted: %v", err)
	}
	want := `{"a":"one","b":"two","c":3}`
	if string(got) != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}
