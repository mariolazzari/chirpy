package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mariolazzari/chirpy/internal/auth"
	"github.com/mariolazzari/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	// check token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing token", err)
		return
	}

	// validate token
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}
	if err = json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// validation
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    userId,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	var rows []database.Chirp
	var err error

	authorID := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	if sortOrder == "" {
		sortOrder = "asc"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		respondWithError(w, http.StatusBadRequest, "invalid sort order", nil)
		return
	}

	if authorID == "" {
		if sortOrder == "desc" {
			rows, err = cfg.db.ReadChirpsDesc(r.Context())
		} else {
			rows, err = cfg.db.ReadChirpsAsc(r.Context())
		}
	} else {
		userID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid user uuid", parseErr)
			return
		}

		if sortOrder == "desc" {
			rows, err = cfg.db.ReadUserChirpsDesc(r.Context(), userID)
		} else {
			rows, err = cfg.db.ReadUserChirpsAsc(r.Context(), userID)
		}
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't read chirps", err)
		return
	}

	chirps := make([]*Chirp, 0, len(rows))
	for _, row := range rows {
		chirps = append(chirps, &Chirp{
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			Body:      row.Body,
			UserID:    row.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	// parse chirp id
	chirpId := r.PathValue("chirpID")
	if chirpId == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID cannot be null", nil)
		return
	}

	// uuid
	chirpUuid, err := uuid.Parse(chirpId)
	if chirpId == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// Read chirp
	row, err := cfg.db.ReadChirp(r.Context(), chirpUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Error searching chirp", err)
		return
	}

	chirp := Chirp{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Body:      row.Body,
		UserID:    row.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// check token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing token", err)
		return
	}

	// validate token
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	// parse chirp id
	chirpId := r.PathValue("chirpID")
	if chirpId == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID cannot be null", nil)
		return
	}

	// uuid
	chirpUuid, err := uuid.Parse(chirpId)
	if chirpId == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// Read chirp
	row, err := cfg.db.ReadChirp(r.Context(), chirpUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Error searching chirp", err)
		return
	}

	// check chirp ownership
	if row.UserID != userId {
		respondWithError(w, http.StatusForbidden, "Error chirp ownership", err)
		return

	}

	err = cfg.db.DeleteChirp(r.Context(), row.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
