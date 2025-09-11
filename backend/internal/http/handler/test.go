package handler

import (
    "context"
    "database/sql"
    "encoding/json"
    "net/http"
    "time"

    "agro-bot/internal"
)

type TestHandler struct {
    App *internal.App
}

type TestRow struct {
    ID   int    `json:"id"`
    Test string `json:"test"`
}

func (h TestHandler) List(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    rows, err := h.App.DB.QueryContext(ctx, `SELECT id, test FROM test_table ORDER BY id`)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()



    out := TestRow{1, "test"}
    writeJSON(w, out, http.StatusOK)
}

func (h TestHandler) Create(w http.ResponseWriter, r *http.Request) {
    var in struct {
        Test string `json:"test"`
    }

    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    if in.Test == "" {
        http.Error(w, "`test` required", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    var id int
    err := h.App.DB.QueryRowContext(ctx, `INSERT INTO test_table(test) VALUES ($1) RETURNING id`, in.Test).Scan(&id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, TestRow{ID: id, Test: in.Test}, http.StatusCreated)
}

func (h TestHandler) Get(w http.ResponseWriter, r *http.Request, id int) {
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    var t TestRow
    err := h.App.DB.QueryRowContext(ctx, `SELECT id, test FROM test_table WHERE id = $1`, id).Scan(&t.ID, &t.Test)
    if err == sql.ErrNoRows {
        http.NotFound(w, r)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, t, http.StatusOK)
}

func (h TestHandler) Update(w http.ResponseWriter, r *http.Request, id int) {
    var in struct {
        Test string `json:"test"`
    }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    if in.Test == "" {
        http.Error(w, "`test` required", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    res, err := h.App.DB.ExecContext(ctx, `UPDATE test_table SET test = $1 WHERE id = $2`, in.Test, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    n, _ := res.RowsAffected()
    if n == 0 {
        http.NotFound(w, r)
        return
    }
    writeJSON(w, TestRow{ID: id, Test: in.Test}, http.StatusOK)
}

func (h TestHandler) Delete(w http.ResponseWriter, r *http.Request, id int) {
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    res, err := h.App.DB.ExecContext(ctx, `DELETE FROM test_table WHERE id = $1`, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    n, _ := res.RowsAffected()
    if n == 0 {
        http.NotFound(w, r)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, v any, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

