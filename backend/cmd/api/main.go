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

    gomavlib "github.com/bluenviron/gomavlib/v2"
    "github.com/bluenviron/gomavlib/v2/pkg/dialects/ardupilotmega"
)

func startMAVLink() {
     node, err := gomavlib.NewNode(gomavlib.NodeConf{
        Endpoints:      []gomavlib.EndpointConf{gomavlib.EndpointUDPServer{Address: "0.0.0.0:14550"}},
        Dialect:        ardupilotmega.Dialect,
        OutVersion:     gomavlib.V2,
        OutSystemID:    255,
        OutComponentID: 190,
    })
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        defer node.Close()
        for evt := range node.Events() {
            if f, ok := evt.(*gomavlib.EventFrame); ok {
                if gp, ok := f.Message().(*ardupilotmega.MessageGlobalPositionInt); ok {
                    ws.DronePosBroadcast(ws.Pos{
                        Lat: float64(gp.Lat) / 1e7,
                        Lon: float64(gp.Lon) / 1e7,
                    })
                }
            }
        }
    }()
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

    startMAVLink()

	mux := http.NewServeMux()

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
