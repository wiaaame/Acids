package repository

import (
	"errors"
	"strings"
	"time"

	"acids/internal/models"

	"gorm.io/gorm"
)

type AcidRepo struct {
	db        *gorm.DB
	minioBase string
}

func NewAcidRepo(db *gorm.DB, minioBase string) *AcidRepo {
	return &AcidRepo{db: db, minioBase: minioBase}
}

func (r *AcidRepo) withBase(a *models.Acid) { a.MinioBase = r.minioBase }

// 1) Первая опубликованная (лента).
func (r *AcidRepo) FirstPublished() (*models.Acid, error) {
	var a models.Acid
	err := r.db.
		Where("status = ?", "published").
		Order("id ASC").
		Preload("Creator").
		Preload("Likes").
		First(&a).Error
	if err != nil {
		return nil, err
	}
	r.withBase(&a)
	return &a, nil
}

// 2) Опубликованная по ID (одна строка через ORM).
func (r *AcidRepo) FindPublishedByID(id uint) (*models.Acid, error) {
	var a models.Acid
	err := r.db.
		Where("id = ? AND status = ?", id, "published").
		Preload("Creator").
		Preload("Likes").
		First(&a).Error
	if err != nil {
		return nil, err
	}
	r.withBase(&a)
	return &a, nil
}

// Следующая опубликованная (циклически).
func (r *AcidRepo) NextPublished(id uint) (*models.Acid, error) {
	var a models.Acid
	err := r.db.
		Where("status = ? AND id > ?", "published", id).
		Order("id ASC").
		Preload("Creator").
		Preload("Likes").
		First(&a).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.FirstPublished()
	}
	if err != nil {
		return nil, err
	}
	r.withBase(&a)
	return &a, nil
}

// 3) Фильтр по молярной массе (ORM, без ручной фильтрации).
func (r *AcidRepo) FilterByMolarMassRange(min, max float64) ([]models.Acid, error) {
	var list []models.Acid
	q := r.db.Where("status = ?", "published")
	if min > 0 {
		q = q.Where("molar_mass >= ?", min)
	}
	if max > 0 {
		q = q.Where("molar_mass <= ?", max)
	}
	if err := q.Order("id ASC").Preload("Creator").Preload("Likes").Find(&list).Error; err != nil {
		return nil, err
	}
	for i := range list {
		r.withBase(&list[i])
	}
	return list, nil
}

// 4) Черновик пользователя.
func (r *AcidRepo) DraftByUser(userID uint) (*models.Acid, error) {
	var a models.Acid
	err := r.db.
		Where("creator_id = ? AND status = ?", userID, "draft").
		Preload("Creator").
		Preload("Likes").
		First(&a).Error
	if err != nil {
		return nil, err
	}
	r.withBase(&a)
	return &a, nil
}

// 5) Создание черновика (кнопка "Далее").
func (r *AcidRepo) CreateDraft(userID uint, reactionName, imageKey, videoKey string) (*models.Acid, error) {
	var existing models.Acid
	err := r.db.Where("creator_id = ? AND status = ?", userID, "draft").First(&existing).Error
	if err == nil {
		existing.ReactionName = reactionName
		existing.ImageKey = imageKey
		existing.VideoKey = videoKey
		existing.UpdatedAt = time.Now()
		if err := r.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		r.withBase(&existing)
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	a := &models.Acid{
		ReactionName: reactionName,
		Status:       "draft",
		ImageKey:     imageKey,
		VideoKey:     videoKey,
		CreatorID:    userID,
	}
	if err := r.db.Create(a).Error; err != nil {
		return nil, err
	}
	r.withBase(a)
	return a, nil
}

// 6) Публикация карточки (кнопка "Опубликовать").
func (r *AcidRepo) Publish(acidID, userID uint, in models.Acid) (*models.Acid, error) {
	var a models.Acid
	err := r.db.
		Where("id = ? AND creator_id = ? AND status = ?", acidID, userID, "draft").
		First(&a).Error
	if err != nil {
		return nil, err
	}

	a.ReactionName = in.ReactionName
	a.AcidName = in.AcidName
	a.AcidFormula = in.AcidFormula
	a.MolarMass = in.MolarMass
	a.Concentration = in.Concentration
	a.CalciumCarbonateMass = in.CalciumCarbonateMass
	a.Description = in.Description
	a.ImageKey = in.ImageKey
	a.VideoKey = in.VideoKey
	a.Status = "published"
	now := time.Now()
	a.PublishedAt = &now
	a.UpdatedAt = now
	if err := r.db.Save(&a).Error; err != nil {
		return nil, err
	}
	r.withBase(&a)
	return &a, nil
}

// ================= Методы веб-сервиса REST (Лабораторная 3) =================

// AcidFilter — параметры фильтрации списка услуг (только опубликованные).
type AcidFilter struct {
	MinMass      float64
	MaxMass      float64
	AcidName     string
	AcidFormula  string
	ReactionName string
}

// ListPublished — список опубликованных услуг с фильтрацией на сервере.
// Записи в статусе deleted на клиент не передаются.
func (r *AcidRepo) ListPublished(f AcidFilter) ([]models.Acid, error) {
	var list []models.Acid
	q := r.db.Where("status = ?", "published")
	if f.MinMass > 0 {
		q = q.Where("molar_mass >= ?", f.MinMass)
	}
	if f.MaxMass > 0 {
		q = q.Where("molar_mass <= ?", f.MaxMass)
	}
	if f.AcidName != "" {
		q = q.Where("LOWER(acid_name) LIKE ?", "%"+strings.ToLower(f.AcidName)+"%")
	}
	if f.AcidFormula != "" {
		q = q.Where("LOWER(acid_formula) LIKE ?", "%"+strings.ToLower(f.AcidFormula)+"%")
	}
	if f.ReactionName != "" {
		q = q.Where("LOWER(reaction_name) LIKE ?", "%"+strings.ToLower(f.ReactionName)+"%")
	}
	if err := q.Order("id ASC").Preload("Creator").Preload("Likes").Find(&list).Error; err != nil {
		return nil, err
	}
	for i := range list {
		r.withBase(&list[i])
	}
	return list, nil
}

// AddOrUpdateDraft — создание услуги (черновик, status=draft).
// Не более одного черновика на пользователя: если черновик уже есть — обновляется.
// Системные поля (статус, создатель, даты) вычисляются на бэкенде, с клиента не принимаются.
func (r *AcidRepo) AddOrUpdateDraft(userID uint, in models.Acid) (*models.Acid, error) {
	var existing models.Acid
	err := r.db.Where("creator_id = ? AND status = ?", userID, "draft").First(&existing).Error
	if err == nil {
		existing.ReactionName = in.ReactionName
		existing.AcidName = in.AcidName
		existing.AcidFormula = in.AcidFormula
		existing.MolarMass = in.MolarMass
		existing.Concentration = in.Concentration
		existing.CalciumCarbonateMass = in.CalciumCarbonateMass
		existing.Description = in.Description
		if in.ImageKey != "" {
			existing.ImageKey = in.ImageKey
		}
		if in.VideoKey != "" {
			existing.VideoKey = in.VideoKey
		}
		existing.UpdatedAt = time.Now()
		if err := r.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		return r.withPreloads(existing.ID)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	a := &models.Acid{
		ReactionName:         in.ReactionName,
		AcidName:             in.AcidName,
		AcidFormula:          in.AcidFormula,
		MolarMass:            in.MolarMass,
		Concentration:        in.Concentration,
		CalciumCarbonateMass: in.CalciumCarbonateMass,
		Description:          in.Description,
		Status:               "draft",
		ImageKey:             in.ImageKey,
		VideoKey:             in.VideoKey,
		CreatorID:            userID,
	}
	if err := r.db.Create(a).Error; err != nil {
		return nil, err
	}
	return r.withPreloads(a.ID)
}

// PublishAcid — смена статуса draft -> published (второй метод создателя).
// Вернуться в черновик нельзя: других переходов статуса API не предоставляет.
func (r *AcidRepo) PublishAcid(acidID, userID uint) (*models.Acid, error) {
	var a models.Acid
	err := r.db.
		Where("id = ? AND creator_id = ? AND status = ?", acidID, userID, "draft").
		First(&a).Error
	if err != nil {
		return nil, err
	}
	now := time.Now()
	a.Status = "published"
	a.PublishedAt = &now
	a.UpdatedAt = now
	if err := r.db.Save(&a).Error; err != nil {
		return nil, err
	}
	return r.withPreloads(a.ID)
}

// SoftDeleteAcid — логическое удаление (soft delete): status = deleted.
// Только услуги владельца; удалённые записи клиенту не возвращаются.
func (r *AcidRepo) SoftDeleteAcid(acidID, userID uint) error {
	var a models.Acid
	err := r.db.
		Where("id = ? AND creator_id = ? AND status <> ?", acidID, userID, "deleted").
		First(&a).Error
	if err != nil {
		return err
	}
	now := time.Now()
	a.Status = "deleted"
	a.CompletedAt = &now
	a.UpdatedAt = now
	return r.db.Save(&a).Error
}

// withPreloads — перечитывает запись с создателем и лайками (для ответов API).
func (r *AcidRepo) withPreloads(id uint) (*models.Acid, error) {
	var a models.Acid
	if err := r.db.Preload("Creator").Preload("Likes").First(&a, id).Error; err != nil {
		return nil, err
	}
	r.withBase(&a)
	return &a, nil
}