package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type PointsHandler struct {
	DB *sql.DB
}

type PointsRow struct {
	ID         int            `json:"id"`
	Latitude   float64        `json:"latitude"`
	Longitude  float64        `json:"longitude"`
	Status     sql.NullString `json:"status"`
	ImageID    sql.NullInt64  `json:"image_id"`
	StartedAt  time.Time      `json:"started_at"`
	FinishedAt sql.NullTime   `json:"finished_at"`
}

func (h PointsHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	rows, err := h.DB.QueryContext(ctx, `SELECT * FROM point ORDER BY id DESC LIMIT 1`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	out := make([]PointsRow, 0, 16)
	for rows.Next() {
		var p PointsRow
		if err := rows.Scan(
			&p.ID,
			&p.Latitude,
			&p.Longitude,
			&p.Status,
			&p.ImageID,
			&p.StartedAt,
			&p.FinishedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		out = append(out, p)
	}
	writeJSON(w, out, http.StatusOK)
}

func (h PointsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	}

	var res struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err := h.DB.QueryRowContext(ctx, `INSERT INTO point(lat, long)
    VALUES ($1, $2)
    RETURNING id`, in.Lat, in.Long).Scan(&res.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, res, http.StatusCreated)
}

// func (h PointsHandler) Get(w http.ResponseWriter, r *http.Request, id int) {
// 	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
// 	defer cancel()

// 	var p PointsRow
// 	err := h.DB.QueryRowContext(ctx, `SELECT * FROM point WHERE id = $1`, id).Scan(&p.id, &p.lat, &p.long, &p.status, &p.image_id, &p.started_at, &p.finished_at)
// 	if err == sql.ErrNoRows {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	writeJSON(w, p, http.StatusOK)
// }

// func (h PointsHandler) Delete(w http.ResponseWriter, r *http.Request, id int) {
// 	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
// 	defer cancel()

// 	res, err := h.DB.ExecContext(ctx, `DELETE FROM point WHERE id = $1`, id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	n, _ := res.RowsAffected()
// 	if n == 0 {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	w.WriteHeader(http.StatusNoContent)
// }
