package router

import (
	"net/http"

	"agro-bot/internal/http/handler"
)

func PointsRouter(mux *http.ServeMux, handler handler.PointsHandler) {
	mux.HandleFunc("/points", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.List(w, r)
		case http.MethodPost:
			handler.Create(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}
