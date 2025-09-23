package router

import (
	"agro-bot/internal/http/handler"
	"net/http"
)

func MissionRouter(mux *http.ServeMux, missionHandler *handler.MissionHandler) {
	mux.HandleFunc("POST /mission", missionHandler.CreateMission)
	mux.HandleFunc("GET /missions", missionHandler.GetMissions)
	mux.HandleFunc("GET /missions/{id}/points", missionHandler.GetMissionPoints)
	mux.HandleFunc("DELETE /missions/{id}", missionHandler.DeleteMission)

}
