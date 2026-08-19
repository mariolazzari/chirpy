package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mariolazzari/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChrp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
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
		UserID: params.UserID,
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
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerReadChirps(w http.ResponseWriter, r *http.Request) {
	rows, err := cfg.db.ReadChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't read chirps", err)
		return
	}

	// create collection
	chirps := make([]Chirp, 0, len(rows))
	for _, row := range rows {
		chirps = append(chirps, Chirp{
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			Body:      row.Body,
			UserID:    row.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerReadChirp(w http.ResponseWriter, r *http.Request) {
	// parse chirp id
	chirpId := r.PathValue("id")
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

// func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
// 	type parameters struct {
// 		Body string `json:"body"`
// 	}
// 	type returnVals struct {
// 		CleanedBody string `json:"cleaned_body"`
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
// 		return
// 	}

// 	const maxChirpLength = 140
// 	if len(params.Body) > maxChirpLength {
// 		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
// 		return
// 	}

// 	badWords := map[string]struct{}{
// 		"kerfuffle": {},
// 		"sharbert":  {},
// 		"fornax":    {},
// 	}
// 	cleaned := getCleanedBody(params.Body, badWords)

// 	respondWithJSON(w, http.StatusOK, returnVals{
// 		CleanedBody: cleaned,
// 	})
// }

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
