package surepay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

// VerifyWebhookSignature verifies the HMAC-SHA256 signature on an inbound
// SurePay webhook event.
//
// secret is your SUREPAY_API_SECRET (same credential used for request signing).
// body is the raw, unparsed request body — do NOT parse it before calling this.
// sig is the expected signature: pass the X-Surepay-Signature header value,
// or leave it empty to extract the signature from body["signature"].
//
// Algorithm (api-flow-sequence.md Flow 9):
//  1. Parse body JSON → map.
//  2. Extract and remove the "signature" field.
//  3. Sort remaining keys alphabetically.
//  4. JSON-encode to compact canonical form (no spaces).
//  5. Compute HMAC-SHA256(secret, canonical_json) → hex.
//  6. Constant-time compare with the extracted signature.
func VerifyWebhookSignature(secret string, body []byte, sig string) bool {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}

	// Resolve signature value
	if sig == "" {
		s, _ := payload["signature"].(string)
		sig = s
	}
	if sig == "" {
		return false
	}
	delete(payload, "signature")

	canonical, err := marshalSorted(payload)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(canonical)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sig))
}

// marshalSorted encodes m to compact JSON with keys sorted alphabetically.
// This produces the canonical form required by the webhook verification algorithm.
func marshalSorted(m map[string]any) ([]byte, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := []byte{'{'}
	for i, k := range keys {
		kb, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		vb, err := json.Marshal(m[k])
		if err != nil {
			return nil, err
		}
		buf = append(buf, kb...)
		buf = append(buf, ':')
		buf = append(buf, vb...)
		if i < len(keys)-1 {
			buf = append(buf, ',')
		}
	}
	buf = append(buf, '}')
	return buf, nil
}
