package handlers

import (
	"html/template"
	"net/http"
	"time"

	"acids/internal/models"
)

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatMM": models.FormatMolarMass,
	}
}

var mediaCheckClient = &http.Client{Timeout: 1500 * time.Millisecond}

// mediaAvailable проверяет доступность файла по HTTP-ссылке (HEAD-запрос).
// Если файл недоступен, шаблон сразу покажет медиа по умолчанию.
func mediaAvailable(url string) bool {
	if url == "" {
		return false
	}
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return false
	}
	resp, err := mediaCheckClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}