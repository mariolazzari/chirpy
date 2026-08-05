package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func MakeRefreshToken() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		// rand.Read should never return an error
		// based on Go's official documentation
		log.Fatal(err)
	}

	refreshToken := hex.EncodeToString(key)

	return refreshToken
}
