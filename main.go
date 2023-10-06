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

const (
	DefaultServerName string = "go-landing"
	ConditionTrue     string = "true"
	ConditionFalse    string = "false"
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

func trimEmpty(slice []string) []string {
	var r []string
	for _, s := range slice {
		if s != "" {
			r = append(r, s)
		}
	}
	return r
}

type GoLandingWriter struct {
	http.ResponseWriter
	Status int
}

func (r *GoLandingWriter) WriteHeader(status int) {
	if r.Status != status {
		r.Status = status
		r.ResponseWriter.WriteHeader(r.Status)
	}
}

func main() {
	serverName := environmentOrDefault("LANDING_SERVER_NAME", DefaultServerName)
	tmpl := environmentOrDefault("LANDING_TEMPLATE", "templates/index.html.tmpl")
	title := environmentOrDefault("LANDING_TITLE", DefaultServerName)
	description := environmentOrDefault("LANDING_DESCRIPTION", "powered by //go:embed")
	links, err := parseLinks(trimEmpty(strings.Split(environmentOrDefault("LANDING_LINKS", ""), ";")))

	catchall := environmentOrDefault("LANDING_CATCHALL_ENABLED", ConditionFalse)
	notFoundMessage := environmentOrDefault("LANDING_NOTFOUND_MESSAGE", "404 Not found")
	ping := environmentOrDefault("LANDING_PING_ENABLED", ConditionTrue)
	pingMessage := environmentOrDefault("LANDING_PING_MESSAGE", "OK")
	hostnames := trimEmpty(strings.Split(environmentOrDefault("LANDING_HOSTNAMES", ""), ";"))

	hostnameMatchingEnabled := (len(hostnames) > 1 || (len(hostnames) == 1 && hostnames[0] != ""))
	if hostnameMatchingEnabled {
		log.Println("hostname matching enabled, allowed hostnames:", hostnames)
	}

	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFS(templateFiles, tmpl)
	if err != nil {
		log.Fatal(err)
	}

	withGoLanding := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			writer := &GoLandingWriter{
				ResponseWriter: w,
				Status:         200,
			}
			writer.Header().Add("X-Server", serverName)

			hostname := strings.ReplaceAll(req.Host, ":9000", "")
			if !strings.HasPrefix(req.URL.Path, "/static/") && hostnameMatchingEnabled {
				match := false
				for _, h := range hostnames {
					if hostname == h {
						match = true
						break
					}
				}
				if !match {
					log.Println("request does not match allowed hostname(s), setting status code to 404")
					writer.WriteHeader(404)
				}
			}

			next.ServeHTTP(writer, req)
			log.Println("serving response", hostname, req.URL.Path, writer.Status)
		})
	}

	pingHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte(pingMessage))
		})
	}

	defaultHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			title := title
			links := links

			var path = req.URL.Path
			if path != "/" && catchall != ConditionTrue {
				w.WriteHeader(404)
				title = notFoundMessage
				links = []Link{{Title: "Back to the startpage", Url: "/"}}
			}
			w.Header().Add("Content-Type", "text/html")

			// respond with the output of template execution
			t.Execute(w, struct {
				Title       string
				Description string
				Links       []Link
			}{Title: title, Description: description, Links: links})
		})
	}

	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	mux := http.NewServeMux()

	if ping == ConditionTrue {
		mux.Handle("/ping", withGoLanding(pingHandler()))
	}
	mux.Handle("/static/", withGoLanding(fs))
	mux.Handle("/", withGoLanding(defaultHandler()))

	log.Println("go-landing is listening on port 9000...")
	err = http.ListenAndServe(":9000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
