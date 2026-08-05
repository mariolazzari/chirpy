package auth

import (
	"encoding/hex"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	expectedLength := 64
	t.Run("successful run", func(t *testing.T) {
		got := MakeRefreshToken()

		assertStringIsNotEmpty(t, got)

		assertTokenHasCorrectLength(t, len(got), expectedLength)

		_, err := hex.DecodeString(got)
		if err != nil {
			t.Errorf("token is not valid hexadecimal: %v", err)
		}
	})

	t.Run("multiple successful runs generate different tokens", func(t *testing.T) {
		token1 := MakeRefreshToken()
		token2 := MakeRefreshToken()

		assertStringsIsNotEqual(t, token1, token2)
	})
}

func assertTokenHasCorrectLength(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("expected token length %d, got %d", want, got)
	}
}
