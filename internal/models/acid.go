package models

import (
	"fmt"
	"time"
)

type Acid struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	ReactionName         string    `gorm:"size:255" json:"reaction_name"`
	AcidName             string    `gorm:"size:100" json:"acid_name"`
	AcidFormula          string    `gorm:"size:50" json:"acid_formula"`
	MolarMass            float64   `json:"molar_mass"`
	Concentration        float64   `json:"concentration"`
	CalciumCarbonateMass float64   `json:"calcium_carbonate_mass"`
	Description          string    `gorm:"type:text" json:"description"`
	Status               string    `gorm:"size:20;not null;default:draft" json:"status"`
	ImageKey             string    `gorm:"size:255;default:''" json:"image_key"`
	VideoKey             string    `gorm:"size:255;default:''" json:"video_key"`
	CreatorID            uint       `gorm:"not null" json:"creator_id"`
	Creator              User       `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	ModeratorID          uint       `gorm:"default:0;not null" json:"moderator_id"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	PublishedAt          *time.Time `json:"published_at"`
	CompletedAt          *time.Time `json:"completed_at"`

	Likes     []Like `gorm:"foreignKey:AcidID" json:"likes,omitempty"`
	MinioBase string `gorm:"-" json:"-"`
}

func (Acid) TableName() string {
	return "acids"
}

// ImageURL - прямая ссылка на изображение в MinIO (или пусто).
func (a *Acid) ImageURL() string {
	if a == nil || a.ImageKey == "" {
		return ""
	}
	return a.MinioBase + a.ImageKey
}

// VideoURL - прямая ссылка на видео в MinIO (или пусто).
func (a *Acid) VideoURL() string {
	if a == nil || a.VideoKey == "" {
		return ""
	}
	return a.MinioBase + a.VideoKey
}

// LikeCount - количество лайков.
func (a *Acid) LikeCount() int {
	if a == nil {
		return 0
	}
	return len(a.Likes)
}

// FormatMolarMass - отображение молярной массы.
func FormatMolarMass(m float64) string { return fmt.Sprintf("%.2f г/моль", m) }