package router

import (
	"agro-bot/internal/http/handler"
	"net/http"
	"strings"
)

func ImageRouter(mux *http.ServeMux, h handler.ImageHandler) {
	mux.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/image/")
		if id == "" || strings.Contains(id, "/") {
			http.NotFound(w, r)
			return
		}
		h.Upload(w, r, id)
	})
}
