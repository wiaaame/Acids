package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"Acids/models"
)

// GridPageData - данные для страницы «Плитка».
type GridPageData struct {
	Acids       []models.Acid
	FilterMin   string // сохранённое значение левого ползунка
	FilterMax   string // сохранённое значение правого ползунка
}

// GridHandler - GET /grid?min={n}&max={m}
// Фильтрует опубликованные услуги по молярной массе кислоты (в диапазоне).
func GridHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		minStr := r.URL.Query().Get("min")
		maxStr := r.URL.Query().Get("max")

		minMass := 0.0
		maxMass := 0.0

		if minStr != "" {
			if v, err := strconv.ParseFloat(minStr, 64); err == nil {
				minMass = v
			}
		}
		if maxStr != "" {
			if v, err := strconv.ParseFloat(maxStr, 64); err == nil {
				maxMass = v
			}
		}

		data := GridPageData{
			Acids:     models.FilterByMolarMassRange(minMass, maxMass),
			FilterMin: minStr,
			FilterMax: maxStr,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "grid.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}