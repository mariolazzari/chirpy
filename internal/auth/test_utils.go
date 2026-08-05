package auth

import "testing"

func assertStringsIsEqual(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertStringsIsNotEqual(t *testing.T, got, want string) {
	if got == want {
		t.Errorf("expected strings to be different")
	}
}

func assertStringIsNotEmpty(t *testing.T, got string) {
	if got == "" {
		t.Fatal("got empty string, expected not empty string")
	}
}
