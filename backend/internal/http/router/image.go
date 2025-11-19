package router

import (
	"agro-bot/internal/http/handler"
	"net/http"
	"path/filepath"
)

func ImageRouter(mux *http.ServeMux, h handler.ImageHandler) {
	mux.HandleFunc("POST /image/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.Upload(w, r, r.PathValue("id"))
	})

	mux.HandleFunc("GET /image", func(w http.ResponseWriter, r *http.Request) {
		h.GetAllImagePaths(w, r)
	})

	mux.HandleFunc("GET /image/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("fileName")
		if fileName == "" {
			http.Error(w, "missing file name", http.StatusBadRequest)
			return
		}

		clean := filepath.Base(fileName)
		path := filepath.Join("uploads", clean)

		http.ServeFile(w, r, path)
	})

}

