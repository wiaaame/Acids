package database

import (
	"database/sql"
	"log"

	"acids/internal/config"
	"acids/internal/models"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	ORM *gorm.DB
	SQL *sql.DB
}

func Connect(cfg config.Config) *DB {
	orm, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения GORM: %v", err)
	}

	sqlDB, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatalf("Ошибка подключения SQL: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ошибка пинга SQL: %v", err)
	}

	// Миграции
	if err := orm.AutoMigrate(&models.User{}, &models.Acid{}, &models.Like{}); err != nil {
		log.Fatalf("Ошибка автомиграции: %v", err)
	}

	return &DB{ORM: orm, SQL: sqlDB}
}