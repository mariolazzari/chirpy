# Learn HTTP Servers by [boot.dev](https://boot.dev)

## http package

[http](https://pkg.go.dev/net/http)

### Server

```go
package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	mux.Handle

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
```

### File server

```go
package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/assets", http.FileServer(http.Dir("./assets")))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe("/healthz", func(w http.re))
}
```

### Handlers

```go
package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
```

## Storage

### Postgres

```sh
docker run --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=chirpy -p 5432:5432 -d postgres:18-alpine
psql "postgres://postgres:@localhost:5432/chirpy"
```

### Goose

```sql
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
```

```sh
goose postgres "postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" up
goose postgres "postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" down
```

### sqlc

[docs](https://sqlc.dev/)

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
```
