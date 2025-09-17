package router

import (
	"agro-bot/internal/http/handler"
	"agro-bot/internal/http/wshandler"
	"net/http"
)

func DroneRouter(mux *http.ServeMux, handler *handler.DroneHandler, handlerWS *wshandler.DroneHandlerWS) {
	mux.HandleFunc("POST /drone/goto", handler.Goto)
	mux.HandleFunc("POST /drone/mission", handler.Mission)
	mux.HandleFunc("DELETE /drone/mission", handler.StopMission)

	mux.HandleFunc("/drone/position", handlerWS.DronePosHandle)
	mux.HandleFunc("/drone/mission/status", handlerWS.DroneMissionHandle)
}
