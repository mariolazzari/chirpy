package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantError bool
	}{
		{
			name:      "valid bearer token",
			header:    "Bearer abc123",
			wantToken: "abc123",
			wantError: false,
		},
		{
			name:      "valid bearer token with extra whitespace",
			header:    "Bearer    abc123",
			wantToken: "abc123",
			wantError: false,
		},
		{
			name:      "missing authorization header",
			header:    "",
			wantToken: "",
			wantError: true,
		},
		{
			name:      "invalid authorization scheme",
			header:    "Basic abc123",
			wantToken: "",
			wantError: true,
		},
		{
			name:      "missing token",
			header:    "Bearer",
			wantToken: "",
			wantError: true,
		},
		{
			name:      "missing token with whitespace",
			header:    "Bearer    ",
			wantToken: "",
			wantError: true,
		},
		{
			name:      "too many parts",
			header:    "Bearer abc123 extra",
			wantToken: "",
			wantError: true,
		},
		{
			name:      "lowercase bearer",
			header:    "bearer abc123",
			wantToken: "abc123",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}

			if tt.header != "" {
				headers.Set("Authorization", tt.header)
			}

			got, err := GetBearerToken(headers)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if got != tt.wantToken {
				t.Errorf("expected token %q, got %q", tt.wantToken, got)
			}
		})
	}
}
