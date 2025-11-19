package handler

import (
	"strconv"
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

type Waypoint struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

type Mission struct {
    ID        int        `json:"id"`
    Waypoints []Waypoint `json:"waypoints"`
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
        SELECT m.id, p.lat, p.long
        FROM mission m
        LEFT JOIN point p ON m.id = p.mission_id
        ORDER BY m.id, p.id
    `)
    if err != nil {
        http.Error(w, "failed to fetch missions", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    missions := make(map[int]*Mission)

    for rows.Next() {
        var (
            missionID int
            lat, lon  *float64
        )

        err := rows.Scan(&missionID, &lat, &lon)
        if err != nil {
            http.Error(w, "scan error", http.StatusInternalServerError)
            return
        }

        if _, exists := missions[missionID]; !exists {
            missions[missionID] = &Mission{
                ID: missionID,
                Waypoints: []Waypoint{},
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

	fmt.Println("Hello?")
	for _, wp := range req.Waypoints {
		fmt.Println(wp)
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
