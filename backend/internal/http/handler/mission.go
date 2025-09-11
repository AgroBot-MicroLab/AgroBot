package handler

import (
	"agro-bot/internal"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type MissionRequest struct {
	Waypoints []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"waypoints"`
}

type MissionHandler struct {
	App *internal.App
}

func (h MissionHandler) CreateMission(w http.ResponseWriter, r *http.Request) {
	var req MissionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := h.App.DB.Begin()
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var missionID int
	err = tx.QueryRow(`INSERT INTO mission DEFAULT VALUES RETURNING id`).Scan(&missionID)
	if err != nil {
		http.Error(w, "failed to create mission", http.StatusInternalServerError)
		return
	}

	for _, wp := range req.Waypoints {
		_, err = tx.Exec(
			`INSERT INTO point (lat, long, mission_id) VALUES ($1, $2, $3)`,
			wp.Lat, wp.Lon, missionID,
		)
		if err != nil {
			http.Error(w, "failed to insert waypoint", http.StatusInternalServerError)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "commit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"mission_id": %d}`, missionID)))
	log.Printf("Created mission %d with %d waypoints", missionID, len(req.Waypoints))
}
