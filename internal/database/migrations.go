package database

import (
	"log"
	"time"

	"acids/internal/models"

	"gorm.io/gorm"
)

// Seed — стартовые данные: пользователи, услуги (draft/published/deleted), лайки.
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

	ts := func() *time.Time { t := time.Now(); return &t }

	acids := []models.Acid{
		{
			// Черновик текущего пользователя (ivanov)
			ReactionName: "Черновик реакции",
			Status:       "draft",
			CreatorID:    users[0].ID,
		},
		{
			// Опубликован (создатель ivanov — совпадает с текущим пользователем)
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
			PublishedAt:          ts(),
		},
		{
			// Удалён — на клиент не передаётся ни одним методом
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
			CompletedAt:          ts(),
		},
		{
			// Опубликован (создатель petrova — is_mine = 0 для текущего пользователя)
			ReactionName:         "CaCO3 + 2HNO3",
			AcidName:             "Азотная кислота",
			AcidFormula:          "HNO3",
			MolarMass:            63.01,
			Concentration:        65,
			CalciumCarbonateMass: 40,
			Description:          "Выделение CO2 при взаимодействии карбоната кальция с азотной кислотой: CaCO3 + 2HNO3 → Ca(NO3)2 + CO2↑ + H2O.",
			Status:               "published",
			ImageKey:             "hno3_acid.jpg",
			VideoKey:             "hno3_reaction.mp4",
			CreatorID:            users[1].ID,
			PublishedAt:          ts(),
		},
		{
			// Опубликован (создатель sidorov)
			ReactionName:         "CaCO3 + 2CH3COOH",
			AcidName:             "Уксусная кислота",
			AcidFormula:          "CH3COOH",
			MolarMass:            60.05,
			Concentration:        70,
			CalciumCarbonateMass: 30,
			Description:          "Растворение карбоната кальция в уксусной кислоте с выделением углекислого газа.",
			Status:               "published",
			ImageKey:             "ch3cooh_acid.jpg",
			VideoKey:             "ch3cooh_reaction.mp4",
			CreatorID:            users[2].ID,
			PublishedAt:          ts(),
		},
	}
	if err := db.Create(&acids).Error; err != nil {
		log.Fatalf("seed acids: %v", err)
	}

	likes := []models.Like{
		{UserID: users[0].ID, AcidID: acids[1].ID},
		{UserID: users[1].ID, AcidID: acids[1].ID},
		{UserID: users[2].ID, AcidID: acids[1].ID},
		{UserID: users[2].ID, AcidID: acids[3].ID},
	}
	if err := db.Create(&likes).Error; err != nil {
		log.Fatalf("seed likes: %v", err)
	}
}