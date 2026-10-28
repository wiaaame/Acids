package models

import "time"

type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_user_acid" json:"user_id"`
	AcidID    uint      `gorm:"not null;uniqueIndex:idx_user_acid" json:"acid_id"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Acid      Acid      `gorm:"foreignKey:AcidID" json:"acid,omitempty"`
}

func (Like) TableName() string {
	return "likes"
}