package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/mariolazzari/chirpy/internal/auth"
)

type loginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresIn int    `json:"expires_in_seconds"`
}

type loginResponse struct {
	User
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid login body", err)
		return
	}

	// search user by email
	user, err := cfg.db.LoginUser(r.Context(), body.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "invalid credentials", nil)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "could not retrieve user", err)
		return
	}

	// verify password
	isValid, err := auth.CheckPasswordHash(body.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not verify password", err)
		return
	}
	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials", nil)
		return
	}

	// token creation
	expiresIn := 3600
	if body.ExpiresIn > 0 {
		expiresIn = min(body.ExpiresIn, 3600)
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, loginResponse{
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: jwt,
	})
}
