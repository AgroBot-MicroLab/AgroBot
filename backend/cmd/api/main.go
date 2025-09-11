package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"agro-bot/internal/db"
	"agro-bot/internal/http/handler"
	"agro-bot/internal/http/middleware"
	"agro-bot/internal/http/router"
	"agro-bot/internal/mav"
	"agro-bot/internal/ws"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	mavc, err := mav.New(mav.Options{
		UDPAddr:         "0.0.0.0:14550",
		OutSystemID:     255,
		OutComponentID:  190,
		TargetSystem:    1,
		TargetComponent: 1,
	})

	if err != nil {
		log.Fatal(err)
	}
	mavc.OnPos = func(p ws.Pos) {

		ws.DronePosBroadcast(ws.Pos{Lat: p.Lat, Lon: p.Lon})

		db.SaveIfChanged(sqlDB, p)
	}

	mavc.OnMissionReached = func(seq uint16) {
		if mavc.MissionActive && mavc.LastSeq == seq {
			log.Println("Mission completed")
			mavc.MissionActive = false
			mavc.LastSeq = 0

			ws.DroneMissionBroadcast(ws.MissionStatus{Status: true})
		}
	}

	mux := http.NewServeMux()

	droneHandler := handler.NewDrone(mavc)
	mux.HandleFunc("/drone/goto", droneHandler.Goto)
	mux.HandleFunc("/drone/mission", droneHandler.Mission)

	testHandler := handler.TestHandler{DB: sqlDB}
	router.TestRouter(mux, testHandler)

	mux.HandleFunc("/drone/position", ws.DronePosHandle)
	mux.HandleFunc("/drone/mission/status", ws.DroneMissionHandle)

	handler := middleware.CorsMiddleware(mux)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

//pointsHandler := handler.PointsHandler{DB: db}
//router.PointsRouter(mux, pointsHandler)

//imgHandler := handler.ImageHandler{DB: db}
//router.ImageRouter(mux, imgHandler)
//router.DeletePointRouter(mux, imgHandler)
