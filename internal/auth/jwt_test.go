package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret-key"
	tokenExpiration := 1 * time.Hour

	t.Run("successful token generation", func(t *testing.T) {
		token, err := MakeJWT(userID, tokenSecret, tokenExpiration)
		if err != nil {
			t.Fatal(err)
		}

		assertTokenIsGenerated(t, token)
	})
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret-key"
	tokenExpiration := 1 * time.Hour

	token, err := MakeJWT(userID, tokenSecret, tokenExpiration)
	if err != nil {
		t.Fatal(err)
	}
	assertTokenIsGenerated(t, token)

	t.Run("verify tokens are equal", func(t *testing.T) {
		userUUID, err := ValidateJWT(token, tokenSecret)
		if err != nil {
			t.Fatal(err)
		}

		assertTokenIsEqual(t, userUUID.String(), userID.String())
	})

	t.Run("validate token with different secret", func(t *testing.T) {
		_, err := ValidateJWT(token, "secret-key-2")
		if err == nil {
			t.Fatal("expected an error but got none")
		}
	})

	t.Run("reject expired token", func(t *testing.T) {
		tokenExpiration = -1 * time.Hour

		expiredToken, err := MakeJWT(userID, tokenSecret, tokenExpiration)
		if err != nil {
			t.Fatal(err)
		}
		assertTokenIsGenerated(t, expiredToken)

		_, err = ValidateJWT(expiredToken, tokenSecret)
		if err == nil {
			t.Fatal("expected error for expired token but got none")
		}
	})

	t.Run("reject malformed token", func(t *testing.T) {
		_, err := ValidateJWT("invalid.token.here", tokenSecret)
		if err == nil {
			t.Fatal("expected error for malformed token but got none")
		}
	})

}

func TestGetBearerToken(t *testing.T) {
	t.Run("sucessful run", func(t *testing.T) {
		got, err := GetBearerToken(http.Header{"Authorization": []string{"Bearer 12341234"}})
		want := "12341234"

		if err != nil {
			t.Fatal(err)
		}

		assertStringsIsEqual(t, got, want)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		_, err := GetBearerToken(http.Header{})

		if err == nil {
			t.Error("expected an error got none")
		}
	})
}

func assertTokenIsGenerated(t *testing.T, token string) {
	if token == "" {
		t.Fatal("generated token is empty")
	}
}

func assertTokenIsEqual(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
