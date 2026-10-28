package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"acids/internal/auth"
	"acids/internal/repository"
)

// DeleteHandler — логическое удаление (soft delete) через ORM.
// Только услуги текущего пользователя (функция-singleton).
func DeleteHandler(repo *repository.AcidRepo) http.HandlerFunc {
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
		id64, err := strconv.ParseUint(parts[3], 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		if err := repo.SoftDeleteAcid(uint(id64), auth.CurrentUser().ID); err != nil {
			http.Error(w, "acid not found or forbidden", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, "/grid", http.StatusSeeOther)
	}
}