package router

import (
	"agro-bot/internal/http/handler"
	"net/http"
	"strings"
)

func ImageRouter(mux *http.ServeMux, h handler.ImageHandler) {
	// Go 1.22+
	mux.HandleFunc("POST /image/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.Upload(w, r, r.PathValue("id"))
	})
	mux.HandleFunc("POST /api/image/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.Upload(w, r, r.PathValue("id"))
	})

	// Совместимо со старым ServeMux (на всякий случай)
	mux.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/image/")
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		h.Upload(w, r, id)
	})
	mux.HandleFunc("/api/image/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/image/")
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		h.Upload(w, r, id)
	})
}
