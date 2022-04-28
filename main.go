package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

type Link struct {
	Title string
	Url   string
}

func parseLinks(slice []string) ([]Link, error) {
	links := make([]Link, 0)
	r := regexp.MustCompile(`\[(?P<Title>.*)\]\((?P<Url>.*)\)`)

	for _, l := range slice {
		if len(l) > 0 {
			p := getParams(r, l)
			links = append(links, Link{Title: p["Title"], Url: p["Url"]})
		}
	}
	return links, nil
}

func getParams(regex *regexp.Regexp, s string) (params map[string]string) {
	match := regex.FindStringSubmatch(s)
	params = make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i > 0 && i <= len(match) {
			params[name] = match[i]
		}
	}
	return params
}

func main() {
	tmpl := environmentOrDefault("LANDING_TEMPLATE", "templates/index.html.tmpl")
	title := environmentOrDefault("LANDING_TITLE", "go-landing")
	description := environmentOrDefault("LANDING_DESCRIPTION", "powered by //go:embed")
	links, err := parseLinks(strings.Split(environmentOrDefault("LANDING_LINKS", ""), ";"))
	catchall := environmentOrDefault("LANDING_CATCHALL", "false")
	notFound := environmentOrDefault("LANDING_NOTFOUND", "404 Not found")

	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFS(templateFiles, tmpl)
	if err != nil {
		log.Fatal(err)
	}

	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	http.Handle("/static/", fs)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		title := title
		links := links

		var path = req.URL.Path
		if path != "/" && catchall != "true" {
			w.WriteHeader(404)
			title = notFound
			links = []Link{{Title: "Back to the startpage", Url: "/"}}
		}

		log.Println("serving request ", path)
		w.Header().Add("Content-Type", "text/html")

		// respond with the output of template execution
		t.Execute(w, struct {
			Title       string
			Description string
			Links       []Link
		}{Title: title, Description: description, Links: links})
	})

	log.Println("go-landing is listening on port 9000...")
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
