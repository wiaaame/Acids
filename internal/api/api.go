package api

import (
	"net/http"

	"acids/internal/auth"
	"acids/internal/repository"
	"acids/internal/storage"
)

// API — REST-веб-сервис (Лабораторная 3) для SPA.
// Реализует два домена: «Услуга» (acids/feed) и «Пользователь» (users).
type API struct {
	acidRepo *repository.AcidRepo
	likeRepo *repository.LikeRepo
	userRepo *repository.UserRepo
	silo     *storage.Silo
}

// New создаёт экземпляр веб-сервиса.
func New(acidRepo *repository.AcidRepo, likeRepo *repository.LikeRepo, userRepo *repository.UserRepo, silo *storage.Silo) *API {
	return &API{acidRepo: acidRepo, likeRepo: likeRepo, userRepo: userRepo, silo: silo}
}

// currentID — функция-singleton: создатель зафиксирован константой в auth.CurrentUser().
// Используется во всех методах веб-сервиса вместо авторизации (до лабораторной 4).
func (a *API) currentID() uint {
	return auth.CurrentUser().ID
}

// RegisterRoutes регистрирует REST-методы с префиксом /api на маршрутизаторе.
func (a *API) RegisterRoutes(mux *http.ServeMux) {
	// Домен «Услуга» (ресурс acids)
	mux.HandleFunc("GET /api/acids", a.listAcids)
	mux.HandleFunc("GET /api/acids/draft", a.getDraft)
	mux.HandleFunc("POST /api/acids", a.createAcid)
	mux.HandleFunc("PUT /api/acids/{id}/publish", a.publishAcid)
	mux.HandleFunc("DELETE /api/acids/{id}", a.deleteAcid)
	mux.HandleFunc("POST /api/acids/{id}/like", a.likeAcid)
	mux.HandleFunc("GET /api/feed", a.feed)

	// Домен «Пользователь»
	mux.HandleFunc("POST /api/users/register", a.register)
	mux.HandleFunc("POST /api/users/login", a.login)
	mux.HandleFunc("POST /api/users/auth", a.login) // алиас: POST аутентификация (формулировка ТЗ)
	mux.HandleFunc("POST /api/users/logout", a.logout)
}