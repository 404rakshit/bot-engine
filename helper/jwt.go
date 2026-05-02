package helper

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a new JWT for a user.
func GenerateToken(userID string, email string, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // 72-hour expiration
		"iat":   time.Now().Unix(),
	})

	return token.SignedString(secret)
}

// VerifyToken parses the token and extracts the core claims (userID and email).
func VerifyToken(tokenString string, secret []byte) (userID string, email string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	// Safely extract the subject (userID)
	userID, ok = claims["sub"].(string)
	if !ok || userID == "" {
		return "", "", errors.New("missing user id in token")
	}

	// Safely extract the email
	email, _ = claims["email"].(string) // We don't fail if email is missing, but we extract it if present

	return userID, email, nil
}
