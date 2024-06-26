package util

import (
	"fmt"
	"time"

	"csv-analyzer-api/internal/entity"

	jwt "github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
}

func CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &JwtCustomClaims{
		UserID: user.ID.String(),

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", fmt.Errorf("invalid Token")
	}

	return claims["id"].(string), nil
}
