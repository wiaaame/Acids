package auth

import (
	"sync"

	"acids/internal/models"
)

// Константы текущего пользователя-создателя (Лабораторная 3).
// Авторизация будет реализована в Лабораторной 4, поэтому во всех методах
// веб-сервиса зафиксирован один и тот же пользователь-создатель.
const (
	CurrentUserID    = uint(1)
	CurrentUsername  = "ivanov"
	CurrentUserEmail = "ivanov@example.com"
)

// Экземпляр singleton'а идёт из функции-конструктора и создаётся один раз (sync.Once),
// поэтому во всех точках приложения используется один и тот же объект пользователя.
var (
	currentUserOnce sync.Once
	currentUser     *models.User
)

// CurrentUser — функция-singleton.
// Всегда возвращает одного и того же пользователя, заданного константами выше.
// Используется во всех методах веб-сервиса вместо авторизации.
func CurrentUser() *models.User {
	currentUserOnce.Do(func() {
		currentUser = &models.User{
			ID:       CurrentUserID,
			Username: CurrentUsername,
			Email:    CurrentUserEmail,
		}
	})
	return currentUser
}