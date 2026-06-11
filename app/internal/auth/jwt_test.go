package auth

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	token, err := GenerateToken(1, 12345)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("token is empty")
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.MasterID != 1 {
		t.Errorf("expected MasterID=1, got %d", claims.MasterID)
	}

	if claims.TelegramID != 12345 {
		t.Errorf("expected TelegramID=12345, got %d", claims.TelegramID)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("invalid-token")
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestValidateToken_Empty(t *testing.T) {
	_, err := ValidateToken("")
	if err != ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken for empty token, got %v", err)
	}
}

func TestTokenClaims_Expiry(t *testing.T) {
	token, err := GenerateToken(1, 12345)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.ExpiresAt == nil {
		t.Fatal("ExpiresAt is nil")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		t.Error("token is already expired")
	}

	futureLimit := time.Now().Add(73 * time.Hour)
	if claims.ExpiresAt.After(futureLimit) {
		t.Error("token expiry is too far in the future")
	}
}
