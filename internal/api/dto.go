package api

import (
	"time"

	"acids/internal/models"
)

// AcidDTO — «сериализатор» услуги для клиента (SPA).
// Системные поля (id, статус, создатель, модератор, даты создания/формирования/завершения)
// формируются только на бэкенде и не принимаются от клиента.
type AcidDTO struct {
	ID                   uint       `json:"id"`
	ReactionName         string     `json:"reaction_name"`
	AcidName             string     `json:"acid_name"`
	AcidFormula          string     `json:"acid_formula"`
	MolarMass            float64    `json:"molar_mass"`
	Concentration        float64    `json:"concentration"`
	CalciumCarbonateMass float64    `json:"calcium_carbonate_mass"`
	Description          string     `json:"description"`
	Status               string     `json:"status"`
	ImageKey             string     `json:"image_key"`
	VideoKey             string     `json:"video_key"`
	ImageURL             string     `json:"image_url"`
	VideoURL             string     `json:"video_url"`
	CreatorName          string     `json:"creator_name"`
	ModeratorID          uint       `json:"moderator_id"`
	IsMine               int        `json:"is_mine"`
	Liked                int        `json:"liked"`
	LikeCount            int        `json:"like_count"`
	CreatedAt            time.Time  `json:"created_at"`
	PublishedAt          *time.Time `json:"published_at"`
	CompletedAt          *time.Time `json:"completed_at"`
}

// toAcidDTO — модель -> DTO.
// is_mine = 1, если создатель услуги совпадает с текущим пользователем (singleton).
// liked   = 1, если текущий пользователь лайкнул услугу.
func toAcidDTO(a *models.Acid, currentUserID uint, liked bool) AcidDTO {
	return AcidDTO{
		ID:                   a.ID,
		ReactionName:         a.ReactionName,
		AcidName:             a.AcidName,
		AcidFormula:          a.AcidFormula,
		MolarMass:            a.MolarMass,
		Concentration:        a.Concentration,
		CalciumCarbonateMass: a.CalciumCarbonateMass,
		Description:          a.Description,
		Status:               a.Status,
		ImageKey:             a.ImageKey,
		VideoKey:             a.VideoKey,
		ImageURL:             a.ImageURL(),
		VideoURL:             a.VideoURL(),
		CreatorName:          a.Creator.Username,
		ModeratorID:          a.ModeratorID,
		IsMine:               boolInt(a.CreatorID == currentUserID),
		Liked:                boolInt(liked),
		LikeCount:            a.LikeCount(),
		CreatedAt:            a.CreatedAt,
		PublishedAt:          a.PublishedAt,
		CompletedAt:          a.CompletedAt,
	}
}

// UserDTO — представление пользователя (пароль никогда не сериализуется).
type UserDTO struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func toUserDTO(u *models.User) UserDTO {
	return UserDTO{ID: u.ID, Username: u.Username, Email: u.Email, CreatedAt: u.CreatedAt}
}