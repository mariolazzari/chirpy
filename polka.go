package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/mariolazzari/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// check api key
	polkaKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing polka key", err)
		return
	}
	if polkaKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid polka key", err)
		return
	}

	// Decode request body
	type data struct {
		UserID string `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	var params parameters

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	// Ignore unsupported events
	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	// Validate user ID
	if params.Data.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing user ID", nil)
		return
	}

	// Upgrade user
	userId, err := uuid.Parse(params.Data.UserID)
	_, err = cfg.db.UpdateUserChirpyRed(r.Context(), uuid.UUID(userId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}
