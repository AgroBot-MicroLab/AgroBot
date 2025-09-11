package db

import (
	"database/sql"
	"log"
	"math"
	"time"

	"agro-bot/internal/ws"
)

func SaveIfChanged(db *sql.DB, pos ws.Pos) {
	const eps = 1e-6

	go func(pos ws.Pos) {
		var lastLat, lastLong float64
		err := db.QueryRow(
			`SELECT lat, "long" FROM drone_position ORDER BY id DESC LIMIT 1`,
		).Scan(&lastLat, &lastLong)

		if err != nil && err != sql.ErrNoRows {
			log.Printf("db error selecting last pos: %v", err)
			return
		}

		needInsert := false
		if err == sql.ErrNoRows {
			needInsert = true
		} else {
			if math.Abs(lastLat-pos.Lat) > eps || math.Abs(lastLong-pos.Lon) > eps {
				needInsert = true
			}
		}

		if needInsert {
			_, err := db.Exec(
				`INSERT INTO drone_position (lat, "long", "timestamp") VALUES ($1, $2, $3)`,
				pos.Lat, pos.Lon, time.Now(),
			)
			if err != nil {
				log.Printf("db insert error: %v", err)
				return
			}
			log.Printf("db: inserted new pos lat=%.6f long=%.6f", pos.Lat, pos.Lon)
		} else {
			log.Printf("db: unchanged pos lat=%.6f long=%.6f — skipped", pos.Lat, pos.Lon)
		}
	}(pos)
}
