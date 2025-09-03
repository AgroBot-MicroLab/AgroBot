package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    _ "github.com/jackc/pgx/v5/stdlib"
)

type resp struct {
    OK      bool   `json:"ok"`
    Version string `json:"version,omitempty"`
    Error   string `json:"error,omitempty"`
}

func main() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL not set")
    }

    db, err := sql.Open("pgx", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := initSchema(db); err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/api/dbcheck", func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()

        var version string
        if err := db.QueryRowContext(ctx, "select version()").Scan(&version); err != nil {
            writeJSON(w, http.StatusServiceUnavailable, resp{OK: false, Error: err.Error()})
            return
        }
        writeJSON(w, http.StatusOK, resp{OK: true, Version: version})
    })

    log.Println("listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func initSchema(db *sql.DB) error {
    _, err := db.Exec(`
        create table if not exists healthprobe(
            id serial primary key,
            created_at timestamptz not null default now()
        );
    `)
    return err
}

func writeJSON(w http.ResponseWriter, code int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(v)
}

