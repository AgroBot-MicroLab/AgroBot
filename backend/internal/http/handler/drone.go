package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"agro-bot/internal"
	"agro-bot/internal/db"
	"agro-bot/internal/mav"
)

type Navigator interface {
	InitMission(ctx context.Context, ptsCount uint16) error
	SendGoto(lat, lon float64) error
	UploadMission(ctx context.Context, wpt []mav.Waypoint) error
	StartMission(ctx context.Context) error
	ClearMissions(ctx context.Context) error
}

type DroneHandler struct {
	App *internal.App
}
type MissionResponse struct {
	MissionID int `json:"mission_id"`
}

func (h *DroneHandler) Goto(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lng"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if err := h.App.MavLinkClient.SendGoto(req.Lat, req.Lon); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *DroneHandler) Mission(w http.ResponseWriter, r *http.Request) {
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

	if err := h.App.MavLinkClient.ClearMissions(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if err := h.App.MavLinkClient.InitMission(ctx, uint16(len(wps))); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if err := h.App.MavLinkClient.UploadMission(ctx, wps); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if err := h.App.MavLinkClient.StartMission(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	missionID, err := db.CreateMission(h.App.DB)
	if err != nil {
		http.Error(w, "failed to create mission", http.StatusInternalServerError)
		return
	}

	resp := MissionResponse{MissionID: missionID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(resp)
}
func (h *DroneHandler) StopMission(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := h.App.MavLinkClient.EndMission(ctx)
	if err != nil {
		log.Printf("failed to stop drone: %v", err)
		http.Error(w, "failed to stop drone", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"stopped"}`))
	log.Println("Drone stopped (LOITER mode)")
}
