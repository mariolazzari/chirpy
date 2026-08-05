package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	params := argon2id.DefaultParams
	return argon2id.CreateHash(password, params)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
