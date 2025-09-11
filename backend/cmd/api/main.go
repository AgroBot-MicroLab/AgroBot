package main

import (
	"agro-bot/internal/http/middleware"
	"agro-bot/internal/mav"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"agro-bot/internal"
	"agro-bot/internal/db"
	"agro-bot/internal/http/handler"
	"agro-bot/internal/http/router"
	"agro-bot/internal/http/wshandler"
	"agro-bot/internal/mqttclient"
	"agro-bot/internal/shared"
)

func main() {
	dbConn := db.NewDBConnection()
	defer dbConn.Close()

	mqttClient := mqttclient.New()
	defer mqttClient.Close()

	mavc, err := mav.New(mav.Options{
		UDPAddr:         "0.0.0.0:14550",
		OutSystemID:     255,
		OutComponentID:  190,
		TargetSystem:    1,
		TargetComponent: 1,
	})

	if err != nil {
		log.Fatalf("mavlink init failed: %v", err)
	}
	defer mavc.Close()

	app := internal.App{DB: dbConn, MavLinkClient: mavc, MqttClient: mqttClient}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /mqtt/test", func(w http.ResponseWriter, r *http.Request) {
		mqttUUID := os.Getenv("MQTT_UUID")
		topic := "agro/" + mqttUUID + "/cmd"
		err := mqttClient.Publish(topic, []byte("make_photo"))
		if err != nil {
			log.Printf("publish error: %v", err)
		}

		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("hello world"))
	})

	testHandler := handler.TestHandler{App: &app}
	router.TestRouter(mux, testHandler)

	imgHandler := handler.ImageHandler{App: &app}
	router.ImageRouter(mux, imgHandler)

	pointsHandler := handler.PointsHandler{App: &app}
	router.PointsRouter(mux, pointsHandler)

	droneHandler := handler.DroneHandler{App: &app}
	droneHandlerWS := wshandler.DroneHandlerWS{App: &app}
	router.DroneRouter(mux, &droneHandler, &droneHandlerWS)

	missionHandler := handler.MissionHandler{App: &app}
	router.MissionRouter(mux, &missionHandler)

	mavc.OnPos = func(p shared.Pos) {
		droneHandlerWS.DronePosBroadcast(shared.Pos{Lat: p.Lat, Lon: p.Lon})
		db.SaveIfChanged(app.DB, p)
	}

	mavc.OnMissionReached = func(seq uint16) {
		if mavc.MissionActive && mavc.LastSeq == seq {
			mavc.MissionActive = false
			mavc.LastSeq = 0
			droneHandlerWS.DroneMissionBroadcast(shared.MissionStatus{Status: true})
		}
	}

	httphandler := middleware.CorsMiddleware(mux)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", httphandler))
}
