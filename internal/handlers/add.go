package handlers

import (
	"html/template"
	"net/http"

	"acids/internal/auth"
	"acids/internal/models"
	"acids/internal/repository"
)

type AddPageData struct {
	HasDraft bool
	Draft    *models.Acid
}

func AddHandler(repo *repository.AcidRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tmpl, err := template.New("").Funcs(templateFuncs()).ParseFiles("templates/add.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		draft, err := repo.DraftByUser(auth.CurrentUser().ID)
		if err == nil && draft != nil {
			// Если файлы недоступны — превью покажет медиа по умолчанию
			if draft.ImageKey != "" && !mediaAvailable(draft.ImageURL()) {
				draft.ImageKey = ""
			}
			if draft.VideoKey != "" && !mediaAvailable(draft.VideoURL()) {
				draft.VideoKey = ""
			}
		}
		data := AddPageData{
			HasDraft: err == nil && draft != nil,
			Draft:    draft,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "add.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}