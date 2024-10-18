// backend/pkg/auth/jwt_test.go
package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	token, err := GenerateJWT(1) // Simulando o usuário com ID 1
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	if token == "" {
		t.Fatalf("Token should not be empty")
	}
}

func TestValidateJWT(t *testing.T) {
	// Gerar um token para validar
	token, err := GenerateJWT(1)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	userID := claims["user_id"].(float64)
	if userID != 1 {
		t.Errorf("Expected user_id 1, but got %v", userID)
	}

	// Testar token expirado
	expiredClaims := jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString(jwtSecretKey)

	_, err = ValidateJWT(expiredTokenString)
	if err == nil {
		t.Errorf("Expected error for expired token, but got none")
	}
}
