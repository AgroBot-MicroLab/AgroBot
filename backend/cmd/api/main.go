package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"agro-bot/internal/http/handler"
	"agro-bot/internal/http/router"
	"agro-bot/internal/mqttclient"
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

	// === MQTT ===
	mqttclient.MustInitFromEnv()
	defer mqttclient.Close()

	// Тестовая публикация при старте
	_ = mqttclient.Publish("AgroBot/test", []byte("hello from backend"))

	// === HTTP-роутеры ===
	mux := http.NewServeMux()

	testHandler := handler.TestHandler{DB: db}
	router.TestRouter(mux, testHandler)

	imgHandler := handler.ImageHandler{DB: db}
	router.ImageRouter(mux, imgHandler)
	router.DeletePointRouter(mux, imgHandler)

	// Health-check на /
	// стало — отдельный health endpoint, не конфликтует
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// === HTTP-сервер ===
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Грейсфул-шутдаун
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	_ = srv.Close()
}
