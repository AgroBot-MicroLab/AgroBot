package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"agro-bot/internal/http/handler"
	"agro-bot/internal/http/router"
	"agro-bot/internal/http/middleware"
	"agro-bot/internal/ws"
    "agro-bot/internal/mav"
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

    mavc, err := mav.New(mav.Options{
        UDPAddr:         "0.0.0.0:14550",
        OutSystemID:     255,
        OutComponentID:  190,
        TargetSystem:    1,
        TargetComponent: 1,
    })

    if err != nil {log.Fatal(err)}
    mavc.OnPos = func(p ws.Pos) {
        ws.DronePosBroadcast(ws.Pos{Lat: p.Lat, Lon: p.Lon});
    }

	mux := http.NewServeMux()

    droneHandler := handler.NewDrone(mavc)
    mux.HandleFunc("/drone/goto", droneHandler.Goto)
    mux.HandleFunc("/drone/mission", droneHandler.Mission)

	testHandler := handler.TestHandler{DB: db}
	router.TestRouter(mux, testHandler)

    mux.HandleFunc("/drone/position", ws.DronePosHandle)

    handler := middleware.CorsMiddleware(mux)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

//pointsHandler := handler.PointsHandler{DB: db}
//router.PointsRouter(mux, pointsHandler)

//imgHandler := handler.ImageHandler{DB: db}
//router.ImageRouter(mux, imgHandler)
//router.DeletePointRouter(mux, imgHandler)
