package jwt

import (
	"fmt"
	"os"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwtlib.RegisteredClaims
}

type JWTManager struct {
	secret string
}

func NewJWTManager() *JWTManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}
	return &JWTManager{secret: secret}
}

func (m *JWTManager) GenerateToken(userID, email string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			Issuer:    "yad-api",
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (m *JWTManager) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwtlib.ParseWithClaims(tokenString, claims, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
