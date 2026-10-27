package handlers

import (
	"html/template"
	"net/http"

	"Acids/models"
)

// AddPageData - данные для страницы «Добавить» (черновик).
type AddPageData struct {
	Draft models.Acid
	All   []models.Acid // для выпадающего списка кислот
}

// AddHandler - GET /add
// Показывает единственную услугу со статусом draft.
func AddHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := AddPageData{
			Draft: models.DraftAcid(),
			All:   models.PublishedAcids(),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "add.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}