package auth

import (
	"net/http"
	"testing"
)

func TestGetApiKey(t *testing.T) {
	t.Run("successfully retrieve API Key", func(t *testing.T) {
		got, err := GetAPIKey(http.Header{"Authorization": []string{"ApiKey l8lD4NQUIt6zUEM+JTNonCBqB/hLuI4P8MNOkfxT04DW9RXuZCm6y2AW6zr9wdhSq8jLOYuflObnE1pZM38yeA=="}})
		want := "l8lD4NQUIt6zUEM+JTNonCBqB/hLuI4P8MNOkfxT04DW9RXuZCm6y2AW6zr9wdhSq8jLOYuflObnE1pZM38yeA=="

		if err != nil {
			t.Fatalf("got an error, expected none: %v", err)
		}

		assertStringsIsEqual(t, got, want)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		_, err := GetAPIKey(http.Header{})

		if err == nil {
			t.Errorf("expected error, got none")
		}

	})
}
