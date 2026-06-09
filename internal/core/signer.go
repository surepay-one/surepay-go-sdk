package core

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// Sign computes the HMAC-SHA256 signature for an outgoing API request.
//
// Algorithm (api-flow-sequence.md Flow 11):
//
//	signing_str = timestamp + "\n" + METHOD + "\n" + path + "\n" + hex(SHA-256(body))
//	signature   = "sha256=" + HMAC-SHA256-hex(secret, signing_str)
func Sign(secret, method, path string, body []byte) (sig, ts string) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	bodyHash := fmt.Sprintf("%x", sha256.Sum256(body))
	signingStr := timestamp + "\n" + method + "\n" + path + "\n" + bodyHash
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingStr))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil)), timestamp
}
