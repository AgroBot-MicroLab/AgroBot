package handler

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

    "agro-bot/internal"
)

const (
	maxUploadBytes = 124 * 1024 * 1024 // 124 Mib
	uploadDir      = "uploads"
)

type ImageHandler struct {
	App *internal.App
}

type ImageRecord struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	SHA256    string    `json:"sha256"`
	CreatedAt time.Time `json:"created_at"`
}

func (h ImageHandler) Upload(w http.ResponseWriter, r *http.Request, pointIDStr string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pointID, err := strconv.ParseInt(pointIDStr, 10, 64)
	if err != nil || pointID <= 0 {
		http.Error(w, "invalid point id", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/form-data") {
		http.Error(w, "expected multipart/form-data", http.StatusBadRequest)
		return
	}
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "missing 'image' field: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	head := make([]byte, 512)
	n, _ := io.ReadFull(file, head)
	ctype := http.DetectContentType(head[:n])
	if !isAllowedImage(ctype) {
		http.Error(w, "unsupported content-type: "+ctype, http.StatusUnsupportedMediaType)
		return
	}
	ext := ".bin"
	if exts, _ := mime.ExtensionsByType(ctype); len(exts) > 0 {
		ext = exts[0]
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := h.App.DB.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var imageID int64
	var capturedAt time.Time

	if err = tx.QueryRowContext(r.Context(),
		`INSERT INTO images (path) VALUES ('') RETURNING id, captured_at`,
	).Scan(&imageID, &capturedAt); err != nil {
		http.Error(w, "db insert images: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%d%s", imageID, ext) // uploads/<imageID>.<ext>
	relPath := filepath.ToSlash(filepath.Join(uploadDir, filename))

	out, e := os.Create(relPath)
	if e != nil {
		err = e
		http.Error(w, "save error: "+e.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	hh := sha256.New()
	size := int64(0)

	if _, e = out.Write(head[:n]); e != nil {
		err = e
		http.Error(w, "write error: "+e.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = hh.Write(head[:n])
	size += int64(n)

	wrote, e := io.Copy(io.MultiWriter(out, hh), file)
	if e != nil {
		err = e
		http.Error(w, "write error: "+e.Error(), http.StatusInternalServerError)
		return
	}
	size += wrote

	if _, e = tx.ExecContext(r.Context(),
		`UPDATE images SET path=$1 WHERE id=$2`, relPath, imageID); e != nil {
		err = e
		http.Error(w, "db update images.path: "+e.Error(), http.StatusInternalServerError)
		return
	}

	res, e := tx.ExecContext(r.Context(),
		`UPDATE point SET image_id=$1 WHERE id=$2 AND (image_id IS DISTINCT FROM $1)`,
		imageID, pointID)
	if e != nil {
		err = e
		_ = os.Remove(relPath)                                                        // чистим файл
		_, _ = tx.ExecContext(r.Context(), `DELETE FROM images WHERE id=$1`, imageID) // чистим запись
		http.Error(w, "update point: "+e.Error(), http.StatusInternalServerError)
		return
	}
	if nrows, _ := res.RowsAffected(); nrows == 0 {
		_ = os.Remove(relPath)
		_, _ = tx.ExecContext(r.Context(), `DELETE FROM images WHERE id=$1`, imageID)
		err = fmt.Errorf("point %d not found", pointID)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if e := tx.Commit(); e != nil {
		err = e
		_ = os.Remove(relPath)
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	err = nil

	resp := map[string]any{
		"image_id":    imageID,
		"point_id":    pointID,
		"path":        relPath,
		"captured_at": capturedAt.UTC(),
		"size":        size,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func isAllowedImage(ct string) bool {
	switch ct {
	case "image/jpeg", "image/png", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}

func sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '-' || r == '_' ||
			(r >= '0' && r <= '9') ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "img"
	}
	return b.String()
}
