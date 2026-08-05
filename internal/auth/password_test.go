package auth

import "testing"

var (
	password = "mypassword123"
)

func TestHashPassword(t *testing.T) {
	t.Run("successful hashing", func(t *testing.T) {
		hashedPassword, err := HashPassword(password)
		if err != nil {
			t.Fatal(err)
		}

		assertStringIsNotEmpty(t, hashedPassword)
		assertStringsIsNotEqual(t, hashedPassword, password)
	})

	t.Run("assert two calls produce different hash", func(t *testing.T) {
		hashedPassword1, err := HashPassword(password)
		if err != nil {
			t.Fatal(err)
		}

		assertStringIsNotEmpty(t, hashedPassword1)

		hashedPassword2, err := HashPassword(password)
		if err != nil {
			t.Fatal(err)
		}

		assertStringIsNotEmpty(t, hashedPassword2)

		assertStringsIsNotEqual(t, hashedPassword1, hashedPassword2)
	})
}

func TestCheckPasswordHash(t *testing.T) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("validate password", func(t *testing.T) {
		isValid, err := CheckPasswordHash(password, hashedPassword)
		if err != nil {
			t.Fatal(err)
		}

		if !isValid {
			t.Errorf("expected matching hashed passwords got not matching hashed passwords")
		}
	})

	t.Run("reject invalid password", func(t *testing.T) {
		wrongHashedPassword, err := HashPassword("wrong-password")
		if err != nil {
			t.Fatal(err)
		}

		isValid, err := CheckPasswordHash(password, wrongHashedPassword)
		if err != nil {
			t.Fatal(err)
		}

		if isValid {
			t.Error("expected password check to fail, but it succeeded")
		}
	})

	t.Run("reject ill formed hash", func(t *testing.T) {
		_, err := CheckPasswordHash(password, "invalid-hash")
		if err == nil {
			t.Errorf("expected error but got none")
		}
	})

}
