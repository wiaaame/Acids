package config

import (
	"log"
	"os"
)

type Config struct {
	DSN       string
	Port      string
	MinioBase string
}

func Load() Config {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=acids port=5432 sslmode=disable TimeZone=UTC"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	minioBase := os.Getenv("MINIO_BASE")
	if minioBase == "" {
		minioBase = "http://localhost:9000/acids/"
	}

	log.Printf("Конфигурация загружена. Порт: %s", port)
	return Config{
		DSN:       dsn,
		Port:      port,
		MinioBase: minioBase,
	}
}