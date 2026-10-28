package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"acids/internal/config"
	"acids/internal/models"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// waitDB — сколько времени сервер ждёт готовности PostgreSQL при старте.
const waitDB = 60 * time.Second

type DB struct {
	ORM *gorm.DB
	SQL *sql.DB
}

// Connect подключается к БД (с ретраями — PostgreSQL может ещё подниматься
// в Docker) и выполняет авто-миграции.
func Connect(cfg config.Config) *DB {
	dsn := withConnectTimeout(cfg.DSN, 5)

	orm, err := connectWithRetry(dsn, waitDB)
	if err != nil {
		log.Fatalf("Ошибка подключения GORM: %v", err)
	}
	log.Printf("Подключение к PostgreSQL установлено")

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения SQL: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ошибка пинга SQL: %v", err)
	}

	// Миграции
	if err := orm.AutoMigrate(&models.User{}, &models.Acid{}, &models.Like{}); err != nil {
		log.Fatalf("Ошибка автомиграции: %v", err)
	}

	return &DB{ORM: orm, SQL: sqlDB}
}

// connectWithRetry пытается подключиться к БД до истечения timeout.
func connectWithRetry(dsn string, timeout time.Duration) (*gorm.DB, error) {
	deadline := time.Now().Add(timeout)
	var lastErr error

	for {
		orm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			if sqlDB, dbErr := orm.DB(); dbErr == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					return orm, nil
				} else {
					err = pingErr
				}
			} else {
				err = dbErr
			}
		}

		lastErr = err
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("БД не готова в течение %s: %w", timeout, lastErr)
		}
		log.Printf("БД ещё не доступна (%v) — повторная попытка...", lastErr)
		time.Sleep(2 * time.Second)
	}
}

// withConnectTimeout добавляет в DSN connect_timeout, если его ещё нет
// (поддерживаются keyword-value и URL-форма postgres://...).
func withConnectTimeout(dsn string, seconds int) string {
	if strings.Contains(dsn, "connect_timeout") {
		return dsn
	}
	suffix := "connect_timeout=" + strconv.Itoa(seconds)
	if strings.Contains(dsn, "://") {
		sep := "?"
		if strings.Contains(dsn, "?") {
			sep = "&"
		}
		return dsn + sep + suffix
	}
	return dsn + " " + suffix
}