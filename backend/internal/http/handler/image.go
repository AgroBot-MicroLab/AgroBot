package handler

import (
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
	"agro-bot/internal/shared"
	"agro-bot/internal/http/wshandler"
)

const (
	maxUploadBytes = 124 * 1024 * 1024 // 124 Mib
	uploadDir      = "uploads"
)

type ImageHandler struct {
	App *internal.App
	DroneHandlerWs *wshandler.DroneHandlerWS
}

type ImageRecord struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	SHA256    string    `json:"sha256"`
	CreatedAt time.Time `json:"created_at"`
}

func (h ImageHandler) GetAllImagePaths(w http.ResponseWriter, r *http.Request) {
    rows, err := h.App.DB.Query(`SELECT path FROM images ORDER BY id`)
    if err != nil {
        http.Error(w, "failed to fetch images", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var paths []string

    for rows.Next() {
        var p string
        if err := rows.Scan(&p); err != nil {
            http.Error(w, "scan error", http.StatusInternalServerError)
            return
        }
        paths = append(paths, p)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, "query error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(paths)
}


// Upload stores an image linked to a mission by mission_id.
// Expects multipart/form-data with field "image".
func (h ImageHandler) Upload(w http.ResponseWriter, r *http.Request, missionIDStr string) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    missionID, err := strconv.ParseInt(missionIDStr, 10, 64)
    if err != nil || missionID <= 0 {
        http.Error(w, "invalid mission id", http.StatusBadRequest)
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

    // Insert shell record with FK to mission. Empty path for now.
    if err = tx.QueryRowContext(
        r.Context(),
        `INSERT INTO images (path, mission_id) VALUES ('', $1) RETURNING id, captured_at`,
        missionID,
    ).Scan(&imageID, &capturedAt); err != nil {
        // If FK fails, likely mission doesn't exist
        http.Error(w, "insert images: "+err.Error(), http.StatusBadRequest)
        return
    }

    filename := fmt.Sprintf("%d%s", imageID, ext)
    relPath := filepath.ToSlash(filepath.Join(uploadDir, filename))

    out, e := os.Create(relPath)
    if e != nil {
        err = e
        http.Error(w, "save error: "+e.Error(), http.StatusInternalServerError)
        return
    }
    defer out.Close()

    size := int64(0)

    if _, e = out.Write(head[:n]); e != nil {
        err = e
        http.Error(w, "write error: "+e.Error(), http.StatusInternalServerError)
        return
    }
    size += int64(n)

    wrote, e := io.Copy(out, file)
    if e != nil {
        err = e
        http.Error(w, "write error: "+e.Error(), http.StatusInternalServerError)
        return
    }
    size += wrote

    if _, e = tx.ExecContext(
        r.Context(),
        `UPDATE images SET path=$1 WHERE id=$2`,
        filename, imageID,
    ); e != nil {
        err = e
        _ = os.Remove(relPath)
        http.Error(w, "db update images.path: "+e.Error(), http.StatusInternalServerError)
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
		"captured_at": capturedAt.UTC(),
	}

	h.DroneHandlerWs.DroneMissionBroadcast(shared.MissionEvent{
		Type: shared.EventPhotoReceived,
		Data: shared.PhotoReceived{
			MissionID: missionID,
			ImageID: imageID,
			Path: filename,
		},
	})

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

