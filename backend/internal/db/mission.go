package db

import (
	"database/sql"
	"fmt"
)

func CreateMission(db *sql.DB) (int64, error) {
    var id int64
    err := db.QueryRow(`INSERT INTO mission DEFAULT VALUES RETURNING id`).Scan(&id)
    if err != nil {
        return 0, err
    }
    fmt.Println("Inserted mission id:", id)
    return id, nil
}

