package handlers

import (
	"net/http"
	"strings"

	"acids/internal/database"
)

// DeleteHandler - логическое удаление через SQL UPDATE (без ORM).
func DeleteHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		idStr := parts[3]

		query := `UPDATE acids SET status = 'deleted', updated_at = NOW() WHERE id = $1 AND status != 'deleted'`
		res, err := db.SQL.Exec(query, idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if n, _ := res.RowsAffected(); n == 0 {
			http.Error(w, "acid not found or already deleted", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/grid", http.StatusSeeOther)
	}
}