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

	// create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// write data to it
	h.Write([]byte(data))

	// get result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil)), nil

}
