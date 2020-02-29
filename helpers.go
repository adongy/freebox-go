package freebox

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

func generatePassword(appToken, challenge string) (string, error) {
	hash := hmac.New(sha1.New, []byte(appToken))
	_, err := hash.Write([]byte(challenge))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
