package router

import (
    "net/http"
    "strconv"
    "strings"

    "agro-bot/internal/http/handler"
)

func TestRouter(mux *http.ServeMux, handler handler.TestHandler) {
    mux.HandleFunc("/tests", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            handler.List(w, r)
        case http.MethodPost:
            handler.Create(w, r)
        default:
            http.NotFound(w, r)
        }
    })

    mux.HandleFunc("/tests/", func(w http.ResponseWriter, r *http.Request) {
        idStr := strings.TrimPrefix(r.URL.Path, "/tests/")
        id, err := strconv.Atoi(idStr)
        if err != nil || id <= 0 {
            http.NotFound(w, r)
            return
        }
        switch r.Method {
        case http.MethodGet:
            handler.Get(w, r, id)
        case http.MethodPut:
            handler.Update(w, r, id)
        case http.MethodDelete:
            handler.Delete(w, r, id)
        default:
            http.NotFound(w, r)
        }
    })
}

