package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHMAC(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyHMAC(data, hmacStr, secret string) bool {
	expected := GenerateHMAC(data, secret)
	return hmacStr == expected
}

func HashCVV(cvv string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(cvv), 12)
	return string(bytes), err
}

func CheckCVV(cvv, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(cvv))
	return err == nil
}
