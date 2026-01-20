package handler

import (
	"agro-bot/internal"
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"crypto/sha256"
    "encoding/hex"
    "sort"
	"strings"
	_ "time"
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

type Waypoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Mission struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	Waypoints   []Waypoint `json:"waypoints"`
}

// PATCH
type MissionUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (h MissionHandler) DeleteMission(w http.ResponseWriter, r *http.Request) {
	missionIDStr := r.PathValue("missionId")
	missionID, err := strconv.Atoi(missionIDStr)
	if err != nil {
		http.Error(w, "invalid mission ID", http.StatusBadRequest)
		return
	}

	tx, err := h.App.DB.Begin()
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM point WHERE mission_id = $1`, missionID)
	if err != nil {
		http.Error(w, "failed to delete mission points", http.StatusInternalServerError)
		return
	}

	res, err := tx.Exec(`DELETE FROM mission WHERE id = $1`, missionID)
	if err != nil {
		http.Error(w, "failed to delete mission", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "mission not found", http.StatusNotFound)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "commit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"deleted"}`))
}

func (h MissionHandler) GetAllMissions(w http.ResponseWriter, r *http.Request) {
	rows, err := h.App.DB.Query(`
        SELECT m.id, m.name, m.description, m.created_at, p.lat, p.long
        FROM mission m
        LEFT JOIN point p ON m.id = p.mission_id
        ORDER BY m.id, p.id
    `)
    if err != nil {
		log.Println("DB query failed:", err)
        http.Error(w, "failed to fetch missions", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

	missions := make(map[int]*Mission)

	for rows.Next() {
		var (
			missionID   int
			name        sql.NullString
			description sql.NullString
			createdAt   time.Time
			lat, lon    *float64
		)

		err := rows.Scan(&missionID, &name, &description, &createdAt, &lat, &lon)
		if err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}

		if _, exists := missions[missionID]; !exists {
			missions[missionID] = &Mission{
				ID:          missionID,
				Name:        name.String,
				Description: description.String,
				CreatedAt:   createdAt,
				Waypoints:   []Waypoint{},
			}
		}

		if lat != nil && lon != nil {
			missions[missionID].Waypoints = append(missions[missionID].Waypoints, Waypoint{
				Lat: *lat,
				Lon: *lon,
			})
		}
	}

	result := make([]Mission, 0, len(missions))
	for _, m := range missions {
		result = append(result, *m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type normWP struct {
    Lat int64
    Lon int64
}

func hashWaypoints(wps []struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}) string {

    const scale = 1e7

    n := make([]normWP, len(wps))
    for i, wp := range wps {
        n[i] = normWP{
            Lat: int64(wp.Lat * scale),
            Lon: int64(wp.Lon * scale),
        }
    }

    sort.Slice(n, func(i, j int) bool {
        if n[i].Lat == n[j].Lat {
            return n[i].Lon < n[j].Lon
        }
        return n[i].Lat < n[j].Lat
    })

    b, _ := json.Marshal(n)
    h := sha256.Sum256(b)
    return hex.EncodeToString(h[:])
}

func (h MissionHandler) CreateMission(w http.ResponseWriter, r *http.Request) {
    var req MissionRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if len(req.Waypoints) == 0 {
        http.Error(w, "no waypoints", http.StatusBadRequest)
        return
    }

    pointsHash := hashWaypoints(req.Waypoints)

    tx, err := h.App.DB.Begin()
    if err != nil {
        http.Error(w, "failed to start transaction", http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    var missionID int
    err = tx.QueryRow(`
        INSERT INTO mission (points_hash)
        VALUES ($1)
        RETURNING id
    `, pointsHash).Scan(&missionID)

    if err != nil {
        if strings.Contains(err.Error(), "mission_points_hash_uq") {
            http.Error(w, "mission with same waypoints already exists", http.StatusConflict)
            return
        }
        http.Error(w, "failed to create mission", http.StatusInternalServerError)
        return
    }

    for _, wp := range req.Waypoints {
        _, err = tx.Exec(
            `INSERT INTO point (lat, long, mission_id)
             VALUES ($1,$2,$3)`,
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
}

func (h MissionHandler) UpdateMission(w http.ResponseWriter, r *http.Request) {
	missionIDStr := r.PathValue("missionId")
	missionID, err := strconv.Atoi(missionIDStr)
	if err != nil {
		http.Error(w, "invalid mission ID", http.StatusBadRequest)
		return
	}

	var req MissionUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == nil && req.Description == nil {
		http.Error(w, "nothing to update", http.StatusBadRequest)
		return
	}

	query := "UPDATE mission SET "
	args := []interface{}{}
	idx := 1

	if req.Name != nil {
		query += fmt.Sprintf("name = $%d", idx)
		args = append(args, *req.Name)
		idx++
	}

	if req.Description != nil {
		if len(args) > 0 {
			query += ", "
		}
		query += fmt.Sprintf("description = $%d", idx)
		args = append(args, *req.Description)
		idx++
	}

	query += fmt.Sprintf(" WHERE id = $%d", idx)
	args = append(args, missionID)

	res, err := h.App.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "failed to update mission", http.StatusInternalServerError)
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "mission not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"updated"}`))
}
