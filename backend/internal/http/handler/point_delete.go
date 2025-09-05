package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Удаляет point по id. Если у точки был image_id, пытаемся удалить images и файл,
// НО только если на этот image_id больше никто не ссылается.
func (h ImageHandler) DeletePoint(w http.ResponseWriter, r *http.Request, idStr string) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pointID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || pointID <= 0 {
		http.Error(w, "invalid point id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1) Удаляем точку и узнаём её image_id (если был)
	var imgID sql.NullInt64
	if err = tx.QueryRowContext(ctx,
		`DELETE FROM point WHERE id=$1 RETURNING image_id`,
		pointID,
	).Scan(&imgID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("point %d not found", pointID), http.StatusNotFound)
			return
		}
		http.Error(w, "delete point: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		imageDeleted bool
		imagePath    sql.NullString
	)

	// 2) Если у точки был image_id — пробуем удалить запись из images,
	//    только если больше нет ссылок из других точек
	if imgID.Valid {
		if err = tx.QueryRowContext(ctx, `
			DELETE FROM images
			WHERE id = $1
			  AND NOT EXISTS (SELECT 1 FROM point WHERE image_id = $1)
			RETURNING path
		`, imgID.Int64).Scan(&imagePath); err != nil {
			if err == sql.ErrNoRows {
				// либо другой point ссылается, либо такого image уже нет — это не ошибка
				err = nil
			} else {
				http.Error(w, "delete image: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			imageDeleted = true
		}
	}

	if e := tx.Commit(); e != nil {
		err = e
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	err = nil

	// 3) Файл удаляем уже после успешного коммита
	if imageDeleted && imagePath.Valid {
		_ = os.Remove(filepath.ToSlash(imagePath.String))
	}

	resp := map[string]any{
		"point_id": pointID,
		"deleted":  true,
		"image_id": func() any {
			if imgID.Valid {
				return imgID.Int64
			}
			return nil
		}(),
		"image_deleted": imageDeleted,
		"image_path": func() any {
			if imagePath.Valid {
				return imagePath.String
			}
			return nil
		}(),
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
