package database

import (
	"log"

	"acids/internal/models"

	"gorm.io/gorm"
)

// Seed - 3 обязательные услуги: draft, published, deleted + пользователи + лайки.
func Seed(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []models.User{
		{Username: "ivanov", Email: "ivanov@example.com", PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeljZAgcfI7p92IdGxad68LJZdL17lhWy"},
		{Username: "petrova", Email: "petrova@example.com", PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeljZAgcfI7p92IdGxad68LJZdL17lhWy"},
		{Username: "sidorov", Email: "sidorov@example.com", PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeljZAgcfI7p92IdGxad68LJZdL17lhWy"},
	}
	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("seed users: %v", err)
	}

	acids := []models.Acid{
		{
			// Черновик
			ReactionName: "Черновик реакции",
			Status:       "draft",
			CreatorID:    users[0].ID,
		},
		{
			// Опубликован
			ReactionName:         "CaCO3 + 2HCl",
			AcidName:             "Соляная кислота",
			AcidFormula:          "HCl",
			MolarMass:            36.46,
			Concentration:        37,
			CalciumCarbonateMass: 50,
			Description:          "Выделение CO2 при взаимодействии карбоната кальция с соляной кислотой. Уравнение: CaCO3 + 2HCl → CaCl2 + CO2↑ + H2O.",
			Status:               "published",
			ImageKey:             "hcl_acid.jpg",
			VideoKey:             "hcl_reaction.mp4",
			CreatorID:            users[0].ID,
		},
		{
			// Удалён
			ReactionName:         "CaCO3 + H2SO4",
			AcidName:             "Серная кислота",
			AcidFormula:          "H2SO4",
			MolarMass:            98.08,
			Concentration:        98,
			CalciumCarbonateMass: 100,
			Description:          "Выделение CO2 при взаимодействии карбоната кальция с серной кислотой.",
			Status:               "deleted",
			ImageKey:             "h2so4_acid.jpg",
			VideoKey:             "h2so4_reaction.mp4",
			CreatorID:            users[1].ID,
		},
	}
	if err := db.Create(&acids).Error; err != nil {
		log.Fatalf("seed acids: %v", err)
	}

	likes := []models.Like{
		{UserID: users[0].ID, AcidID: acids[1].ID},
		{UserID: users[1].ID, AcidID: acids[1].ID},
		{UserID: users[2].ID, AcidID: acids[1].ID},
	}
	if err := db.Create(&likes).Error; err != nil {
		log.Fatalf("seed likes: %v", err)
	}
}