package main

import (
	"log"
	"net/http"

	"acids/internal/api"
	"acids/internal/config"
	"acids/internal/database"
	"acids/internal/handlers"
	"acids/internal/repository"
	"acids/internal/storage"
)

func main() {
	cfg := config.Load()

	// Подключение к БД и миграции
	db := database.Connect(cfg)
	database.Seed(db.ORM)

	// Репозитории (взаимодействие с БД только через ORM)
	acidRepo := repository.NewAcidRepo(db.ORM, cfg.MinioBase)
	likeRepo := repository.NewLikeRepo(db.ORM)
	userRepo := repository.NewUserRepo(db.ORM)

	// Объектное хранилище Silo (MinIO): корзина + политика публичного чтения
	silo, err := storage.NewSilo(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioBase,
		cfg.MinioUseSSL,
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации Silo/MinIO: %v", err)
	}

	mux := http.NewServeMux()

	// ---------- Лабораторные 1-2: HTML-страницы ----------
	mux.HandleFunc("GET /feed", handlers.FeedHandler(acidRepo))
	mux.HandleFunc("GET /add", handlers.AddHandler(acidRepo))
	mux.HandleFunc("POST /create", handlers.CreateDraftHandler(acidRepo))
	mux.HandleFunc("POST /publish", handlers.PublishHandler(acidRepo))
	mux.HandleFunc("GET /grid", handlers.GridHandler(acidRepo))
	mux.HandleFunc("POST /grid/delete/", handlers.DeleteHandler(acidRepo))

	// ---------- Лабораторная 3: REST-веб-сервис /api ----------
	apiSvc := api.New(acidRepo, likeRepo, userRepo, silo)
	apiSvc.RegisterRoutes(mux)

	// Статика
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	log.Printf("Сервер запущен на порту %s (REST API: /api, Silo: %s)", cfg.Port, cfg.MinioBase)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}