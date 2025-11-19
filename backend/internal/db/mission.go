package db

import (
	"agro-bot/internal/mav"
	"database/sql"

	"fmt"
)

func CreateMission(db *sql.DB, wps []mav.Waypoint) (int64, error) {
    var id int64
    err := db.QueryRow(`INSERT INTO mission DEFAULT VALUES RETURNING id`).Scan(&id)
    if err != nil {
        return 0, err
    }

	for _, wp := range wps {
		fmt.Println(wp)
		_, err = db.Exec(
			`INSERT INTO point (lat, long, mission_id) VALUES ($1, $2, $3)`,
			wp.Lat, wp.Lon, id,
		)

		if err != nil {
			return 0, err
		}
	}

    fmt.Println("Inserted mission id:", id)
    return id, nil
}

