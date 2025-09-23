package handler

import (
	"agro-bot/internal"
	"encoding/json"
	//	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Waypoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type MissionRequest struct {
	Waypoints []Waypoint `json:"waypoints"`
}

type MissionHandler struct {
	App *internal.App
}

func (h MissionHandler) CreateMission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req MissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if len(req.Waypoints) == 0 {
		http.Error(w, `{"error":"waypoints required"}`, http.StatusBadRequest)
		return
	}
	for _, wp := range req.Waypoints {
		if wp.Lat < -90 || wp.Lat > 90 || wp.Lon < -180 || wp.Lon > 180 {
			http.Error(w, `{"error":"invalid coordinates"}`, http.StatusBadRequest)
			return
		}
	}

	tx, err := h.App.DB.Begin()
	if err != nil {
		http.Error(w, `{"error":"failed to start transaction"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var missionID int
	if err := tx.QueryRow(`INSERT INTO mission DEFAULT VALUES RETURNING id`).Scan(&missionID); err != nil {
		http.Error(w, `{"error":"failed to create mission"}`, http.StatusInternalServerError)
		return
	}

	for _, wp := range req.Waypoints {
		if _, err := tx.Exec(
			`INSERT INTO point (lat, long, mission_id) VALUES ($1, $2, $3)`,
			wp.Lat, wp.Lon, missionID,
		); err != nil {
			http.Error(w, `{"error":"failed to insert waypoint"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, `{"error":"commit failed"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"mission_id": missionID})
	log.Printf("Created mission %d with %d waypoints", missionID, len(req.Waypoints))
}

func (h MissionHandler) GetMissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := h.App.DB.Query(`SELECT id, created_at FROM mission ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch missions"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Mission struct {
		ID        int       `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}

	var missions []Mission
	for rows.Next() {
		var m Mission
		if err := rows.Scan(&m.ID, &m.CreatedAt); err != nil {
			http.Error(w, `{"error":"failed to scan mission"}`, http.StatusInternalServerError)
			return
		}
		missions = append(missions, m)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, `{"error":"rows error"}`, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(missions)
}

func (h MissionHandler) GetMissionPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id") // Go 1.22+
	missionID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid mission id"}`, http.StatusBadRequest)
		return
	}

	rows, err := h.App.DB.Query(
		`SELECT lat, long FROM point WHERE mission_id = $1 ORDER BY id`, missionID,
	)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch points"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type WP struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}
	var out []WP

	for rows.Next() {
		var p WP
		if err := rows.Scan(&p.Lat, &p.Lon); err != nil {
			http.Error(w, `{"error":"scan error"}`, http.StatusInternalServerError)
			return
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, `{"error":"rows error"}`, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(out)
}

func (h MissionHandler) DeleteMission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid mission id"}`, http.StatusBadRequest)
		return
	}
	if _, err := h.App.DB.Exec(`DELETE FROM mission WHERE id = $1`, id); err != nil {
		http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
