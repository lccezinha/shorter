package main

import (
  "fmt"
  "log"
  "net/http"
  "strings"
  "encoding/json"

  "github.com/lccezinha/shorter/url"
)

type Headers map[string]string

type Redirecter struct {
  stats chan string
}

var (
  port int
  urlBase string
)

func init() {
  port = 8080
  urlBase = fmt.Sprintf("http://localhost:%d", port)
}

func respondWith(w http.ResponseWriter, status int, headers Headers) {
  for key, value := range headers {
    w.Header().Set(key, value)
  }

  w.WriteHeader(status)
}

func respondWithJSON(w http.ResponseWriter, responseData string) {
  respondWith(w, http.StatusOK, Headers{"Content-type":"application/json"})
  fmt.Fprintf(w, responseData)
}

func logger(format string, values ...interface{}) {
  log.Printf(fmt.Sprintf("%s \n", format), values...)
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

func (redirecter *Redirecter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  path := strings.Split(r.URL.Path, "/")
  id := path[len(path) - 1]

  if url := url.Find(id); url != nil {
    http.Redirect(w, r, url.UrlOriginal, http.StatusMovedPermanently)
    redirecter.stats <- id
  } else {
    http.NotFound(w, r)
  }
}

func Shorter(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    respondWith(w, http.StatusMethodNotAllowed, Headers{"Allow":"POST"})
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
    "Link": fmt.Sprintf("<%s/api/stats/%s>; rel='stats'", urlBase, url.Id),
  }
  respondWith(w, status, headers)

  logger("URL: %s foi encurtada para %s. \n", url.UrlOriginal, urlShort)
}

func Stats(w http.ResponseWriter, r *http.Request) {
  path := strings.Split(r.URL.Path, "/")
  id := path[len(path) - 1]

  logger("Buscando clicks do ID: %s", id)

  if url := url.Find(id); url != nil {
    json, err := json.Marshal(url.ShowStats())
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    respondWithJSON(w, string(json))
  } else {
    http.NotFound(w, r)
  }
}

func Home(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello")
}

func main() {
  stats := make(chan string)
  defer close(stats)
  go registerStats(stats)

  url.ConfigRepository(url.InitializeRepository())

  http.HandleFunc("/api/shorter", Shorter)
  http.Handle("/r/", &Redirecter{stats})
  http.HandleFunc("/home", Home)
  http.HandleFunc("/api/stats/", Stats)

  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}