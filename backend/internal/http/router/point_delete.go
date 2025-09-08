package router

import (
	"net/http"
	"strings"

	"agro-bot/internal/http/handler"
)

func DeletePointRouter(mux *http.ServeMux, h handler.ImageHandler) {
	// Новый стиль (Go 1.22+)
	mux.HandleFunc("DELETE /delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.DeletePoint(w, r, r.PathValue("id"))
	})
	mux.HandleFunc("DELETE /delite/{id}", func(w http.ResponseWriter, r *http.Request) { // по твоему названию
		h.DeletePoint(w, r, r.PathValue("id"))
	})

	// Совместимо со старым ServeMux
	mux.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.NotFound(w, r)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/delete/")
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		h.DeletePoint(w, r, id)
	})
	mux.HandleFunc("/delite/", func(w http.ResponseWriter, r *http.Request) { // по твоему названию
		if r.Method != http.MethodDelete {
			http.NotFound(w, r)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/delite/")
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		h.DeletePoint(w, r, id)
	})
}
