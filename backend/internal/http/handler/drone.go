package handler

import (
    "encoding/json"
    "net/http"
    "time"
    "context"

    "agro-bot/internal/mav"
)

type Navigator interface {
    InitMission(ctx context.Context, ptsCount uint16) error
    SendGoto(lat, lon float64) error
    UploadMission(ctx context.Context, wpt []mav.Waypoint) error
    StartMission(ctx context.Context) error
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

    wps := []mav.Waypoint{
        {Lat: -35.36214764686344, Lon: 149.1651090448245},
        {Lat: -35.36214764686344, Lon: 149.1661090448245},
        {Lat: -35.36264000000000, Lon: 149.1666000000000},
        {Lat: -35.36264000000000, Lon: 149.1671000000000},
    }

    if err := h.Nav.InitMission(ctx, uint16(len(wps))); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }

    if err := h.Nav.UploadMission(ctx, wps); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }

    if err := h.Nav.StartMission(ctx); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }

    w.WriteHeader(http.StatusOK)
}


