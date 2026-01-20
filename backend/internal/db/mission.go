package db

import (
	"agro-bot/internal/mav"
	"database/sql"

	"crypto/sha256"
	"encoding/hex"
	"sort"
	"encoding/json"
	"errors"
)


func hashWaypoints(wps []mav.Waypoint) string {
    type norm struct {
        Lat int64
        Lon int64
    }

    const scale = 1e7

    n := make([]norm, len(wps)-1)
	for i, wp := range wps[1:] {
        n[i] = norm{
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


func CreateMission(db *sql.DB, wps []mav.Waypoint) (int64, error) {
    if len(wps) == 0 {
        return 0, errors.New("no waypoints")
    }

    pointsHash := hashWaypoints(wps)

    var id int64
    err := db.QueryRow(
        `INSERT INTO mission (points_hash)
         VALUES ($1)
         RETURNING id`,
        pointsHash,
    ).Scan(&id)

    if err != nil {
        return 0, err
    }

	for _, wp := range wps[1:] {
        _, err = db.Exec(
            `INSERT INTO point (lat, long, mission_id)
             VALUES ($1,$2,$3)`,
            wp.Lat, wp.Lon, id,
        )
        if err != nil {
            return 0, err
        }
    }

    return id, nil
}

