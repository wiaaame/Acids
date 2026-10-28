package api

import (
	"net/http"
	"strings"

	"acids/internal/auth"
)

// RegisterRequest — тело запроса регистрации.
// Системные поля (ид, даты) с клиента не принимаются.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// register — POST /api/users/register.
// Регистрация нового пользователя (пароль хешируется bcrypt).
func (a *API) register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := readJSON(r, &req); err != nil {
		writeErr(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		writeErr(w, http.StatusBadRequest, "username, email и password обязательны")
		return
	}

	u, err := a.userRepo.Create(req.Username, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already") {
			writeErr(w, http.StatusConflict, err.Error())
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toUserDTO(u))
}

// login — POST /api/users/login.
// Аутентификация: заглушка для лабораторной 4.
func (a *API) login(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"message":      "Аутентификация — заглушка, реализуется в лабораторной 4",
		"session":      "stub-session",
		"current_user": toUserDTO(auth.CurrentUser()),
	})
}

// logout — POST /api/users/logout.
// Деавторизация: заглушка для лабораторной 4.
func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Деавторизация — заглушка, реализуется в лабораторной 4",
	})
}