package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sync/atomic"
	"unicode/utf8"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {

	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, apiCfg.fileserverHits.Load())
	})

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits.Store(0)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hits: %d", apiCfg.fileserverHits.Load())
	})

	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		type Body struct {
			Body string `json:"body"`
		}

		type BodyError struct {
			Message string `json:"message"`
		}

		type Response struct {
			CleanedBody string `json:"cleaned_body"`
		}

		w.Header().Set("Content-Type", "application/json")

		var body Body
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(BodyError{
				Message: "Something went wrong",
			})
			return
		}

		// validation
		msgLen := utf8.RuneCountInString(body.Body)
		if msgLen < 1 || msgLen > 140 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(BodyError{
				Message: "Chirp is too long",
			})
			return
		}

		// cleaned
		re := regexp.MustCompile(`(?i)\b(kerfuffle|sharbert|fornax)\b`)
		cleaned := re.ReplaceAllString(body.Body, "****")

		json.NewEncoder(w).Encode(Response{
			CleanedBody: cleaned,
		})
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
