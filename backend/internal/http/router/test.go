package router

import (
    "net/http"
    "strconv"
    "strings"

    "agro-bot/internal/http/handler"
)

type Handlers struct {
    Test handler.TestHandler
}

func TestRouter(mux *http.ServeMux, h Handlers) {
    mux.HandleFunc("/tests", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            h.Test.List(w, r)
        case http.MethodPost:
            h.Test.Create(w, r)
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
            h.Test.Get(w, r, id)
        case http.MethodPut:
            h.Test.Update(w, r, id)
        case http.MethodDelete:
            h.Test.Delete(w, r, id)
        default:
            http.NotFound(w, r)
        }
    })
}

