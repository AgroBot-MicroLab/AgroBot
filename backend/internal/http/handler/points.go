package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

    "agro-bot/internal"
)

type PointsHandler struct {
    App *internal.App
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

	rows, err := h.App.DB.QueryContext(ctx, `SELECT * FROM point ORDER BY id DESC LIMIT 1`)
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

	err := h.App.DB.QueryRowContext(ctx, `INSERT INTO point(lat, long)
    VALUES ($1, $2)
    RETURNING id`, in.Lat, in.Long).Scan(&res.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, res, http.StatusCreated)
}

