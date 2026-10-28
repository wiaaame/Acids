package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"acids/internal/models"
	"acids/internal/repository"
)

type FeedPageData struct {
	Acid    *models.Acid
	NextID  uint
	HasMore bool
}

func FeedHandler(repo *repository.AcidRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tmpl, err := template.New("").Funcs(templateFuncs()).ParseFiles("templates/feed.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		q := r.URL.Query()
		idStr := q.Get("id")
		nextStr := q.Get("next")

		var current *models.Acid

		if idStr == "" {
			current, err = repo.FirstPublished()
			if err != nil {
				http.Error(w, "no published acids", http.StatusNotFound)
				return
			}
		} else {
			id, convErr := strconv.ParseUint(idStr, 10, 64)
			if convErr != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}

			if nextStr == "true" {
				current, err = repo.NextPublished(uint(id))
			} else {
				current, err = repo.FindPublishedByID(uint(id))
			}
			if err != nil {
				http.Error(w, "acid not found", http.StatusNotFound)
				return
			}
		}

		// Обрезка описания для ленты (прямо в указателе — безопасно, т.к. current != nil)
		if len([]rune(current.Description)) > 120 {
			current.Description = string([]rune(current.Description)[:120]) + "..."
		}

		// Если файлы недоступны по ссылке — шаблон сразу отобразит медиа по умолчанию
		if current.ImageKey != "" && !mediaAvailable(current.ImageURL()) {
			current.ImageKey = ""
		}
		if current.VideoKey != "" && !mediaAvailable(current.VideoURL()) {
			current.VideoKey = ""
		}

		var nextID uint
		if next, err := repo.NextPublished(current.ID); err == nil {
			nextID = next.ID
		}

		data := FeedPageData{Acid: current, NextID: nextID, HasMore: true}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "feed.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}