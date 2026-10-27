package main

import (
	"log"
	"net/http"

	"acids/internal/config"
	"acids/internal/database"
	"acids/internal/handlers"
	"acids/internal/repository"
)

func main() {
	cfg := config.Load()

	// Подключение к БД и миграции
	db := database.Connect(cfg)
	database.Seed(db.ORM)

	// Репозитории
	acidRepo := repository.NewAcidRepo(db.ORM, cfg.MinioBase)

	// Маршруты
	http.HandleFunc("/feed", handlers.FeedHandler(acidRepo))
	http.HandleFunc("/add", handlers.AddHandler(acidRepo))
	http.HandleFunc("/create", handlers.CreateDraftHandler(acidRepo))
	http.HandleFunc("/publish", handlers.PublishHandler(acidRepo))
	http.HandleFunc("/grid", handlers.GridHandler(acidRepo))
	http.HandleFunc("/grid/delete/", handlers.DeleteHandler(db))

	// Статика
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Printf("Сервер запущен на порту %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}