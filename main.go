package main

import (
	"html/template"
	"log"
	"net/http"

	"Acids/handlers"
)

func main() {
	// Загружаем все шаблоны.
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("template parse error: %v", err)
	}

	mux := http.NewServeMux()

	// Три GET-контроллера по предметной области.
	mux.HandleFunc("/feed", handlers.FeedHandler(tmpl))   // лента по ID + next=true
	mux.HandleFunc("/add", handlers.AddHandler(tmpl))     // единственный черновик
	mux.HandleFunc("/grid", handlers.GridHandler(tmpl))   // плитка + фильтр по молярной массе

	// Статика (CSS).
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Корень - редирект на /feed.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/feed", http.StatusFound)
	})

	addr := ":3001"
	log.Printf("server started at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}