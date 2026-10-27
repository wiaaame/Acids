package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"Acids/models"
)

// FeedPageData - данные, передаваемые в шаблон feed.html
type FeedPageData struct {
	Acid    models.Acid
	NextID  int
	HasMore bool // показывать ли кнопку «Ещё»
}

// FeedHandler - GET /feed?id={id}&next=true
// Без параметров - открывает первую опубликованную.
// ?id=N - открывает услугу с указанным ID.
// ?id=N&next=true - открывает следующую после N.
func FeedHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		q := r.URL.Query()
		idStr := q.Get("id")
		nextStr := q.Get("next")

		var current models.Acid
		var ok bool

		if idStr == "" {
			// без параметров – первая опубликованная
			firstID := models.FirstPublishedID()
			current, ok = models.FindByID(firstID)
			if !ok {
				http.Error(w, "no published acids", http.StatusNotFound)
				return
			}
		} else {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}
			if nextStr == "true" {
				current = models.NextPublished(id)
			} else {
				current, ok = models.FindByID(id)
				if !ok {
					http.Error(w, "acid not found", http.StatusNotFound)
					return
				}
			}
		}

		// Обрезаем описание до 2 строк (условно — первые 120 символов)
		if len([]rune(current.Description)) > 120 {
			current.Description = string([]rune(current.Description)[:120]) + "…"
		}

		data := FeedPageData{
			Acid:    current,
			NextID:  models.NextPublished(current.ID).ID,
			HasMore: true,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "feed.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}