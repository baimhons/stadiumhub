package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/baimhons/stadiumhub/internal"
)

func GenerateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

var secretKey = internal.ENV.SecretKey.SecretKey

// SignSession สร้างลายเซ็น HMAC จาก sessionID
func SignSession(sessionID string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(sessionID))
	signature := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(signature)
}

// VerifySession ตรวจสอบว่าลายเซ็นถูกต้องไหม
func VerifySession(sessionID, signature string) bool {
	expected := SignSession(sessionID)
	return hmac.Equal([]byte(expected), []byte(signature))
}
