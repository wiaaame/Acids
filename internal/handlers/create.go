package handlers

import (
	"net/http"

	"acids/internal/auth"
	"acids/internal/repository"
)

func CreateDraftHandler(repo *repository.AcidRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		reactionName := r.FormValue("reaction_name")
		imageKey := r.FormValue("image_key")
		videoKey := r.FormValue("video_key")

		if _, err := repo.CreateDraft(auth.CurrentUser().ID, reactionName, imageKey, videoKey); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/add", http.StatusSeeOther)
	}
}