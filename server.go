package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lccezinha/shorter/url"
)

type Headers map[string]string

type Redirecter struct {
	stats chan string
}

var (
	port    *int
	logOn   *bool
	urlBase string
	logFile *os.File
)

func init() {
	port = flag.Int("p", 8888, "port")
	logOn = flag.Bool("l", true, "log on/off")
	flag.Parse()

	urlBase = fmt.Sprintf("http://localhost:%d", *port)
	logFile, _ = os.Create("server.log")
}

func respondWith(w http.ResponseWriter, status int, headers Headers) {
	for key, value := range headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(status)
}

func respondWithJSON(w http.ResponseWriter, responseData string) {
	respondWith(w, http.StatusOK, Headers{"Content-type": "application/json"})
	fmt.Fprintf(w, responseData)
}

func logger(format string, values ...interface{}) {
	if *logOn {
		log.Printf(fmt.Sprintf("%s \n", format), values...)
		logFile.WriteString(fmt.Sprintf(fmt.Sprintf("%s \n", format), values...))
	}
}

func extractUrl(r *http.Request) string {
	url := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(url)
	return string(url)
}

func registerStats(ids <-chan string) {
	for id := range ids {
		url.RegisterClick(id)
		logger("Click registrado para %s. ", id)
	}
}

func fetchAndExecute(w http.ResponseWriter, r *http.Request, executor func(*url.Url)) {
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

	if url := url.Find(id); url != nil {
		executor(url)
	} else {
		http.NotFound(w, r)
	}
}

func (redirecter *Redirecter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fetchAndExecute(w, r, func(url *url.Url) {
		http.Redirect(w, r, url.UrlOriginal, http.StatusMovedPermanently)
		redirecter.stats <- url.Id
		logger("Fazendo redirect da URL %s \n", url.UrlOriginal)
	})
}

func Shorter(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respondWith(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})
		return
	}

	url, isNew, err := url.FetchUrl(extractUrl(r))

	if err != nil {
		respondWith(w, http.StatusBadRequest, nil)
		return
	}

	var status int
	if isNew {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	urlShort := fmt.Sprintf("%s/r/%s", urlBase, url.Id)
	headers := Headers{
		"Location": urlShort,
		"Link":     fmt.Sprintf("<%s/api/stats/%s>; rel='stats'", urlBase, url.Id),
	}
	respondWith(w, status, headers)

	logger("URL: %s foi encurtada para %s. \n", url.UrlOriginal, urlShort)
}

func Stats(w http.ResponseWriter, r *http.Request) {
	fetchAndExecute(w, r, func(url *url.Url) {
		json, err := json.Marshal(url.ShowStats())

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, string(json))
	})
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func main() {
	stats := make(chan string)
	defer close(stats)
	go registerStats(stats)
	defer logFile.Close()

	url.ConfigRepository(url.InitializeMongoRepository())

	http.HandleFunc("/api/shorter", Shorter)
	http.Handle("/r/", &Redirecter{stats})
	http.HandleFunc("/home", Home)
	http.HandleFunc("/api/stats/", Stats)

	logger("Server iniciado na porta: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
