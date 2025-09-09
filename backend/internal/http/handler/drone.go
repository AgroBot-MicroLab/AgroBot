package handler

import (
    "encoding/json"
    "net/http"
    "time"
    "context"
)

type Navigator interface {
    SendGoto(lat, lon float64) error
    RunHardcodedMission(ctx context.Context) (error)
    StartMission(ctx context.Context) (error)
}

type Drone struct {
    Nav Navigator
}


func NewDrone(nav Navigator) *Drone {
    return &Drone{Nav: nav}
}

func (h *Drone) Goto(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Lat float64 `json:"lat"`
        Lon float64 `json:"lon"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad json", http.StatusBadRequest); return
    }

    if err := h.Nav.SendGoto(req.Lat, req.Lon); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway); return
    }

    w.WriteHeader(http.StatusAccepted)
}

func (h *Drone) Mission(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := h.Nav.RunHardcodedMission(ctx); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }

    if err := h.Nav.StartMission(ctx); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }

    w.WriteHeader(http.StatusOK)
}


