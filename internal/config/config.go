package config

import (
	"log"
	"os"
)

type Config struct {
	DSN       string
	Port      string
	MinioBase string

	// Silo (MinIO) — объектное хранилище для файлов услуг
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func Load() Config {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=acids port=5432 sslmode=disable TimeZone=UTC connect_timeout=5"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	minioBase := os.Getenv("MINIO_BASE")
	if minioBase == "" {
		minioBase = "http://localhost:9000/acids/"
	}

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "localhost:9000"
	}

	minioAccess := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccess == "" {
		minioAccess = "minioadmin"
	}

	minioSecret := os.Getenv("MINIO_SECRET_KEY")
	if minioSecret == "" {
		minioSecret = "minioadmin"
	}

	minioBucket := os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		minioBucket = "acids"
	}

	minioUseSSL := os.Getenv("MINIO_USE_SSL") == "true"

	log.Printf("Конфигурация загружена. Порт: %s, Silo: %s, bucket: %s", port, minioEndpoint, minioBucket)
	return Config{
		DSN:            dsn,
		Port:           port,
		MinioBase:      minioBase,
		MinioEndpoint:  minioEndpoint,
		MinioAccessKey: minioAccess,
		MinioSecretKey: minioSecret,
		MinioBucket:    minioBucket,
		MinioUseSSL:    minioUseSSL,
	}
}