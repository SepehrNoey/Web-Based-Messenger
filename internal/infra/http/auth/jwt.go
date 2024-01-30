package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	secretKey          []byte
	expirationDuration time.Duration
	signMethod         jwt.SigningMethod
	iss                string
	aud                []string
}

func NewJWTConfig(secretKey []byte, expDur time.Duration, signMethod jwt.SigningMethod, iss string, aud []string) *JWTConfig {
	return &JWTConfig{
		secretKey:          secretKey,
		expirationDuration: expDur,
		signMethod:         signMethod,
		iss:                iss,
		aud:                aud,
	}
}

func (c *JWTConfig) CreateToken(claims map[string]interface{}) (string, error) {
	allClaims := jwt.MapClaims{}
	allClaims["iss"] = c.iss
	allClaims["aud"] = c.aud
	for key, value := range claims {
		allClaims[key] = value
	}

	token := jwt.NewWithClaims(c.signMethod, allClaims)

	tokenString, err := token.SignedString(c.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c *JWTConfig) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if c.signMethod != t.Method {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return c.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}

	return claims, nil
}
