package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
)

func environmentOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

//go:embed templates
var templateFiles embed.FS

//go:embed static
var staticFiles embed.FS

func main() {
	tmpl := environmentOrDefault("LANDING_TEMPLATE", "templates/index.html.tmpl")
	title := environmentOrDefault("LANDING_TITLE", "go-landing")
	description := environmentOrDefault("LANDING_DESCRIPTION", "powered by //go:embed")

	t, err := template.ParseFS(templateFiles, tmpl)
	if err != nil {
		log.Fatal(err)
	}

	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	http.Handle("/static/", fs)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var path = req.URL.Path
		if path != "/" {
			w.WriteHeader(404)
			w.Write([]byte(""))
			return
		}

		log.Println("serving request ", path)
		w.Header().Add("Content-Type", "text/html")

		// respond with the output of template execution
		t.Execute(w, struct {
			Title       string
			Description string
		}{Title: title, Description: description})
	})

	log.Println("go-landing is listening on port 9000...")
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
