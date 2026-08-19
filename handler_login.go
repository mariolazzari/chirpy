package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mariolazzari/chirpy/internal/auth"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid login body", err)
		return
	}

	user, err := cfg.db.LoginUser(r.Context(), body.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "invalid credentials", nil)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "could not retrieve user", err)
		return
	}

	isValid, err := auth.CheckPasswordHash(body.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not verify password", err)
		return
	}

	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
