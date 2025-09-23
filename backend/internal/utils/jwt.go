package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT interface {
	Generate(claimsMap map[string]interface{}, expUnix int64, secret string) (string, error)
	Parse(tokenString string, claims jwt.Claims, secret string) (*jwt.Token, error)
}

type jwtService struct{}

func NewJWT() JWT {
	return &jwtService{}
}

// Generate สร้าง token พร้อม claim
func (j *jwtService) Generate(claimsMap map[string]interface{}, expUnix int64, secret string) (string, error) {
	claims := jwt.MapClaims{}
	for k, v := range claimsMap {
		claims[k] = v
	}
	claims["exp"] = expUnix
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Parse ตรวจสอบความถูกต้องของ token และดึง claims ออกมา
func (j *jwtService) Parse(tokenString string, claims jwt.Claims, secret string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่าใช้ signing method ที่ถูกต้อง
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
