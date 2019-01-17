package checkout

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// GenerateHMAC generates HMAC SHA-256 hash
func GenerateHMAC(secret string, data string) (string, error) {

	if secret == "" || data == "" {
		return "", errors.New("secret or data is missing")
	}

	// create a new HMAC SHA-256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))

	// return result as hexadecimal string
	return hex.EncodeToString(h.Sum(nil)), nil

}
