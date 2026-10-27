package repository

import (
	"errors"

	"acids/internal/models"

	"gorm.io/gorm"
)

type LikeRepo struct {
	db *gorm.DB
}

func NewLikeRepo(db *gorm.DB) *LikeRepo {
	return &LikeRepo{db: db}
}

func (r *LikeRepo) CountByAcid(acidID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("acid_id = ?", acidID).Count(&count).Error
	return count, err
}

func (r *LikeRepo) CountByAcids(acidIDs []uint) (map[uint]int64, error) {
	result := make(map[uint]int64)
	if len(acidIDs) == 0 {
		return result, nil
	}
	type row struct {
		AcidID uint
		Total  int64
	}

	var rows []row
	err := r.db.Model(&models.Like{}).
		Select("acid_id, COUNT(*) AS total").
		Where("acid_id IN ?", acidIDs).
		Group("acid_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r.AcidID] = r.Total
	}
	return result, nil
}

func (r *LikeRepo) Exists(userID, acidID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Like{}).
		Where("user_id = ? AND acid_id = ?", userID, acidID).
		Count(&count).Error
	return count > 0, err
}

func (r *LikeRepo) Create(userID, acidID uint) (*models.Like, error) {
	var acid models.Acid
	if err := r.db.Where("id = ? AND status = ?", acidID, "published").First(&acid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("acid not found or not published")
		}
		return nil, err
	}
	exists, err := r.Exists(userID, acidID)
	if err != nil {
		return nil, err
	}
	if exists {
		var existing models.Like
		if err := r.db.Where("user_id = ? AND acid_id = ?", userID, acidID).First(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}
	like := &models.Like{UserID: userID, AcidID: acidID}
	if err := r.db.Create(like).Error; err != nil {
		return nil, err
	}
	return like, nil
}

func (r *LikeRepo) Delete(userID, acidID uint) error {
	return r.db.Where("user_id = ? AND acid_id = ?", userID, acidID).Delete(&models.Like{}).Error
}

func (r *LikeRepo) Toggle(userID, acidID uint) (bool, int64, error) {
	exists, err := r.Exists(userID, acidID)
	if err != nil {
		return false, 0, err
	}
	if exists {
		if err := r.Delete(userID, acidID); err != nil {
			return false, 0, err
		}
	} else {
		if _, err := r.Create(userID, acidID); err != nil {
			return false, 0, err
		}
	}
	count, err := r.CountByAcid(acidID)
	if err != nil {
		return !exists, 0, err
	}
	return !exists, count, nil
}

func (r *LikeRepo) ListByAcid(acidID uint) ([]models.Like, error) {
	var likes []models.Like
	err := r.db.Where("acid_id = ?", acidID).Preload("User").Order("created_at DESC").Find(&likes).Error
	return likes, err
}

func (r *LikeRepo) ListByUser(userID uint) ([]models.Like, error) {
	var likes []models.Like
	err := r.db.Where("user_id = ?", userID).Preload("Acid").Order("created_at DESC").Find(&likes).Error
	return likes, err
}

func (r *LikeRepo) DeleteAllByAcid(acidID uint) error {
	return r.db.Where("acid_id = ?", acidID).Delete(&models.Like{}).Error
}