package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"acids/internal/repository"
)

type AcidView struct {
	ID           uint
	ReactionName string
	AcidName     string
	AcidFormula  string
	MolarMass    float64
	ImageURL     string
	LikeCount    int
}

type GridPageData struct {
	Acids     []AcidView
	FilterMin string
	FilterMax string
}

func GridHandler(repo *repository.AcidRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tmpl, err := template.New("").Funcs(templateFuncs()).ParseFiles("templates/grid.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		minStr := r.URL.Query().Get("min")
		maxStr := r.URL.Query().Get("max")

		var minMass, maxMass float64
		if v, err := strconv.ParseFloat(minStr, 64); err == nil {
			minMass = v
		}
		if v, err := strconv.ParseFloat(maxStr, 64); err == nil {
			maxMass = v
		}

		list, err := repo.FilterByMolarMassRange(minMass, maxMass)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		views := make([]AcidView, 0, len(list))
		for _, a := range list {
			imageURL := a.ImageURL()
			// Если файл недоступен — шаблон сразу покажет фото по умолчанию
			if imageURL != "" && !mediaAvailable(imageURL) {
				imageURL = ""
			}
			views = append(views, AcidView{
				ID:           a.ID,
				ReactionName: a.ReactionName,
				AcidName:     a.AcidName,
				AcidFormula:  a.AcidFormula,
				MolarMass:    a.MolarMass,
				ImageURL:     imageURL,
				LikeCount:    a.LikeCount(),
			})
		}

		data := GridPageData{Acids: views, FilterMin: minStr, FilterMax: maxStr}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "grid.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}