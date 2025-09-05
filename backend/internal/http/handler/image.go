package handler

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxUploadBytes = 20 << 20 // 20 MB
	uploadDir      = "uploads"
)

type ImageHandler struct {
	DB *sql.DB
}

type ImageRecord struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	SHA256    string    `json:"sha256"`
	CreatedAt time.Time `json:"created_at"`
}

func (h ImageHandler) Upload(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if id == "" || strings.Contains(id, "/") {
		http.Error(w, "invalid id", http.StatusBadRequest)
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
	detected := http.DetectContentType(head[:n])
	if !isAllowedImage(detected) {
		http.Error(w, "unsupported content-type: "+detected, http.StatusUnsupportedMediaType)
		return
	}

	ext := ".bin"
	if exts, _ := mime.ExtensionsByType(detected); len(exts) > 0 {
		ext = exts[0]
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := fmt.Sprintf("%s_%d%s", sanitize(id), time.Now().UnixNano(), ext)
	path := filepath.Join(uploadDir, filename)

	out, err := os.Create(path)
	if err != nil {
		http.Error(w, "save error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	hh := sha256.New()
	var size int64

	if _, err := out.Write(head[:n]); err != nil {
		http.Error(w, "write error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = hh.Write(head[:n])
	size += int64(n)

	wrote, err := io.Copy(io.MultiWriter(out, hh), file)
	if err != nil {
		var me *multipart.Part
		if errors.As(err, &me) {
			http.Error(w, "multipart error: "+err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "write error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	size += wrote

	// В БД пишем только path (как в твоей схеме), id и captured_at получаем через RETURNING
	storedPath := filepath.ToSlash(path) // чтобы в БД были прямые слэши
	var dbID int64
	var capturedAt time.Time

	if err := h.DB.QueryRow(
		`INSERT INTO images (path) VALUES ($1) RETURNING id, captured_at`,
		storedPath,
	).Scan(&dbID, &capturedAt); err != nil {
		_ = os.Remove(path)
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ответ клиенту (ничего в схеме не меняем)
	resp := map[string]any{
		"id":          dbID,       // PK из БД
		"path":        storedPath, // что записали в БД
		"captured_at": capturedAt.UTC(),
		"drone_id":    id,   // id из URL, для информации
		"size":        size, // сервисная инфа (не в БД)
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
