package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;unique;not null" json:"username"`
	Email        string    `gorm:"size:100;unique;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Acids        []Acid    `gorm:"foreignKey:CreatorID" json:"acids,omitempty"`
	Likes        []Like    `gorm:"foreignKey:UserID" json:"likes,omitempty"`
}

func (User) TableName() string {
	return "users"
}