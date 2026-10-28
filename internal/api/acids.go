package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"acids/internal/models"
	"acids/internal/repository"

	"gorm.io/gorm"
)

// maxBodySize — максимальный размер multipart-запроса (поля + файлы), 60 МБ.
const maxBodySize = 60 << 20

// parseID — извлечение ид услуги из пути /api/acids/{id}/...
func parseID(r *http.Request) (uint, error) {
	v, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

// ---------- GET /api/acids ----------

// listAcids — список услуг (ресурс acids, только опубликованные) с фильтрацией на бэкенде.
// Параметры: min, max (молярная масса), acid_name, acid_formula, reaction_name.
// Признак is_mine = 1, если создатель услуги совпадает с текущим пользователем.
func (a *API) listAcids(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := repository.AcidFilter{
		MinMass:      parseFloat(q.Get("min")),
		MaxMass:      parseFloat(q.Get("max")),
		AcidName:     strings.TrimSpace(q.Get("acid_name")),
		AcidFormula:  strings.TrimSpace(q.Get("acid_formula")),
		ReactionName: strings.TrimSpace(q.Get("reaction_name")),
	}

	list, err := a.acidRepo.ListPublished(f)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "ошибка получения списка: "+err.Error())
		return
	}

	// Признак liked вычисляется массово (один запрос), is_mine — сравнением с создателем.
	uid := a.currentID()
	ids := make([]uint, 0, len(list))
	for i := range list {
		ids = append(ids, list[i].ID)
	}
	likedMap, err := a.likeRepo.MapByUser(uid, ids)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "ошибка подсчёта лайков: "+err.Error())
		return
	}

	dtos := make([]AcidDTO, 0, len(list))
	for i := range list {
		dtos = append(dtos, toAcidDTO(&list[i], uid, likedMap[list[i].ID]))
	}
	writeJSON(w, http.StatusOK, dtos)
}

// ---------- GET /api/feed ----------

// feed — лента опубликованных услуг (одна запись в ответе):
// без id — первая опубликованная, ?id=N — по ид, ?id=N&next=true — следующая (циклически).
// Признак liked = 1, если текущий пользователь лайкнул эту услугу.
func (a *API) feed(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	idStr := q.Get("id")
	uid := a.currentID()

	var (
		acid *models.Acid
		err  error
	)
	switch {
	case idStr == "":
		acid, err = a.acidRepo.FirstPublished()
	case q.Get("next") == "true":
		id, conv := strconv.ParseUint(idStr, 10, 64)
		if conv != nil {
			writeErr(w, http.StatusBadRequest, "некорректный id")
			return
		}
		acid, err = a.acidRepo.NextPublished(uint(id))
	default:
		id, conv := strconv.ParseUint(idStr, 10, 64)
		if conv != nil {
			writeErr(w, http.StatusBadRequest, "некорректный id")
			return
		}
		acid, err = a.acidRepo.FindPublishedByID(uint(id))
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeErr(w, http.StatusNotFound, "услуга не найдена")
			return
		}
		writeErr(w, http.StatusInternalServerError, "ошибка загрузки ленты: "+err.Error())
		return
	}

	liked, err := a.likeRepo.Exists(uid, acid.ID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "ошибка проверки лайка: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toAcidDTO(acid, uid, liked))
}

// ---------- GET /api/acids/draft ----------

// getDraft — черновик текущего пользователя (не более одной записи), ид не указывается.
func (a *API) getDraft(w http.ResponseWriter, r *http.Request) {
	uid := a.currentID()
	draft, err := a.acidRepo.DraftByUser(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeErr(w, http.StatusNotFound, "черновик не найден")
			return
		}
		writeErr(w, http.StatusInternalServerError, "ошибка загрузки черновика: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toAcidDTO(draft, uid, false))
}

// ---------- POST /api/acids ----------

// createAcid — добавление услуги (multipart/form-data) с файлами изображения и видео.
// Файлы загружаются в Silo (MinIO) с латинскими именами, в БД сохраняются ключи.
// Системные поля (ид, статус, создатель, даты) на бэкенде, с клиента не принимаются.
func (a *API) createAcid(w http.ResponseWriter, r *http.Request) {
	uid := a.currentID()

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "too large") {
			status = http.StatusRequestEntityTooLarge
		}
		writeErr(w, status, "не удалось разобрать форму: "+err.Error())
		return
	}

	in := models.Acid{
		ReactionName: strings.TrimSpace(r.FormValue("reaction_name")),
		AcidName:     strings.TrimSpace(r.FormValue("acid_name")),
		AcidFormula:  strings.TrimSpace(r.FormValue("acid_formula")),
		Description:  strings.TrimSpace(r.FormValue("description")),
	}
	in.MolarMass = parseFloat(r.FormValue("molar_mass"))
	in.Concentration = parseFloat(r.FormValue("concentration"))
	in.CalciumCarbonateMass = parseFloat(r.FormValue("calcium_carbonate_mass"))

	// Файлы как файлы, а НЕ URL: image — изображение, video — короткое видео.
	var imageKey, videoKey string
	if parts := r.MultipartForm.File["image"]; len(parts) > 0 {
		key, err := a.silo.UploadImage(parts[0])
		if err != nil {
			writeErr(w, http.StatusBadRequest, "изображение: "+err.Error())
			return
		}
		imageKey = key
	}
	if parts := r.MultipartForm.File["video"]; len(parts) > 0 {
		key, err := a.silo.UploadVideo(parts[0])
		if err != nil {
			writeErr(w, http.StatusBadRequest, "видео: "+err.Error())
			return
		}
		videoKey = key
	}
	in.ImageKey = imageKey
	in.VideoKey = videoKey

	if strings.TrimSpace(in.ReactionName) == "" && imageKey == "" && videoKey == "" {
		writeErr(w, http.StatusBadRequest, "заполните поля услуги или прикрепите файлы")
		return
	}

	saved, err := a.acidRepo.AddOrUpdateDraft(uid, in)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "ошибка сохранения: "+err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toAcidDTO(saved, uid, false))
}

// ---------- PUT /api/acids/{id}/publish ----------

// publishAcid — смена статуса draft -> published (второй метод создателя).
// Вернуться в черновик нельзя: других переходов статуса API не предоставляет.
func (a *API) publishAcid(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "некорректный id")
		return
	}
	uid := a.currentID()

	saved, err := a.acidRepo.PublishAcid(id, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeErr(w, http.StatusNotFound, "черновик не найден или недоступен")
			return
		}
		writeErr(w, http.StatusInternalServerError, "ошибка публикации: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toAcidDTO(saved, uid, false))
}

// ---------- DELETE /api/acids/{id} ----------

// deleteAcid — логическое удаление услуги (soft delete, status = deleted).
// Только услуги текущего пользователя; удалённые записи клиенту не возвращаются.
func (a *API) deleteAcid(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "некорректный id")
		return
	}
	if err := a.acidRepo.SoftDeleteAcid(id, a.currentID()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeErr(w, http.StatusNotFound, "услуга не найдена или недоступна")
			return
		}
		writeErr(w, http.StatusInternalServerError, "ошибка удаления: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "услуга удалена", "id": id, "deleted": 1})
}

// ---------- POST /api/acids/{id}/like ----------

// likeAcid — лайк от текущего пользователя.
// Тело: {"value": 1} — поставить лайк, {"value": 0} — снять лайк.
func (a *API) likeAcid(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "некорректный id")
		return
	}

	var req struct {
		Value *int `json:"value"`
	}
	if err := readJSON(r, &req); err != nil {
		writeErr(w, http.StatusBadRequest, "ожидается JSON вида {\"value\": 0|1}")
		return
	}
	if req.Value == nil || (*req.Value != 0 && *req.Value != 1) {
		writeErr(w, http.StatusBadRequest, "поле value должно быть 0 или 1")
		return
	}

	liked, count, err := a.likeRepo.SetLike(a.currentID(), id, *req.Value == 1)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not published") {
			writeErr(w, http.StatusNotFound, err.Error())
			return
		}
		writeErr(w, http.StatusInternalServerError, "ошибка лайка: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"acid_id":    id,
		"liked":      boolInt(liked),
		"like_count": count,
	})
}

// parseFloat — безопасный разбор числа из строки запроса/формы.
func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}