package handlers

import (
	"net/http"
	"strconv"

	"acids/internal/models"
	"acids/internal/repository"
)

func PublishHandler(repo *repository.AcidRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		acidID64, _ := strconv.ParseUint(r.FormValue("acid_id"), 10, 64)
		molar, _ := strconv.ParseFloat(r.FormValue("molar_mass"), 64)
		conc, _ := strconv.ParseFloat(r.FormValue("concentration"), 64)

		// Название и сведения о кислоте заполняются на шаге «Далее» и хранятся в черновике;
		// форма публикации содержит только описание и параметры, поэтому недостающие поля берём из черновика.
		draft, err := repo.DraftByUser(1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if draft.ID != uint(acidID64) {
			http.Error(w, "draft not found", http.StatusNotFound)
			return
		}

		in := models.Acid{
			ReactionName:  draft.ReactionName,
			AcidName:      draft.AcidName,
			AcidFormula:   draft.AcidFormula,
			MolarMass:     molar,
			Concentration: conc,
			Description:   r.FormValue("description"),
			ImageKey:      r.FormValue("image_key"),
			VideoKey:      r.FormValue("video_key"),
		}

		if _, err := repo.Publish(uint(acidID64), 1, in); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/feed", http.StatusSeeOther)
	}
}