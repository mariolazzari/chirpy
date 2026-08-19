package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	secret := "my-secret"

	t.Run("create and validate", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT() error = %v", err)
		}

		gotUserID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("ValidateJWT() error = %v", err)
		}

		if gotUserID != userID {
			t.Errorf("ValidateJWT() userID = %v, want %v", gotUserID, userID)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, -time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT() error = %v", err)
		}

		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Error("ValidateJWT() expected error for expired token")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("MakeJWT() error = %v", err)
		}

		_, err = ValidateJWT(token, "wrong-secret")
		if err == nil {
			t.Error("ValidateJWT() expected error for wrong secret")
		}
	})
}
