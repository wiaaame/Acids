package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSON — сериализация ответа в JSON c нужным HTTP-статусом.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("api: ошибка записи JSON: %v", err)
	}
}

// writeErr — ответ с ошибкой вида {"error": "..."}.
func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// readJSON — разбор тела запроса (неизвестные поля запрещены — защита системных полей).
func readJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// boolInt переводит bool в признак 0/1 (требуется для лайка и признаков is_mine/liked).
func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}