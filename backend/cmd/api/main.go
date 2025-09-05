package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    _ "github.com/jackc/pgx/v5/stdlib"

    "agro-bot/internal/http/handler"
    "agro-bot/internal/http/router"
)

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

    mux := http.NewServeMux()

    testHandler := handler.TestHandler{DB: db}
    router.TestRouter(mux, testHandler);
    
    pointsHandler := handler.PointsHandler{DB: db}
    router.PointsRouter(mux, pointsHandler)

    log.Println("listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

