package router

import (
	"agro-bot/internal/http/handler"

	"net/http"
)

func MissionRouter(mux *http.ServeMux, missionHandler *handler.MissionHandler) {
	mux.HandleFunc("POST /mission", func(w http.ResponseWriter, r *http.Request) {
		missionHandler.CreateMission(w, r)
	})

}
