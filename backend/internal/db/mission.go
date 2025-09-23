package db

import (
	"database/sql"
	"log"
)

func CreateMission(db *sql.DB) (int, error) {
	var missionID int
	err := db.QueryRow(
		`INSERT INTO mission DEFAULT VALUES RETURNING id`,
	).Scan(&missionID)
	if err != nil {
		log.Printf("db: mission creation error: %v", err)
		return 0, err
	}
	log.Printf("db: mission created with id=%d", missionID)
	return missionID, nil
}
