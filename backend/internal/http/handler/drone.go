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
    ClearMissions(ctx context.Context) error
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
        Lon float64 `json:"lng"`
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

    var wps []mav.Waypoint
    if err := json.NewDecoder(r.Body).Decode(&wps); err != nil {
        http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    if len(wps) == 0 {
        http.Error(w, "no waypoints provided", http.StatusBadRequest)
        return
    }

    if err := h.Nav.ClearMissions(ctx); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
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


