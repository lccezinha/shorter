package main

import (
  "fmt"
  "log"
  "net/http"
  "strings"

  "github.com/lccezinha/shorter/url"
)

type Headers map[string]string

var (
  port int
  urlBase string
)

func init() {
  port = 8080
  urlBase = fmt.Sprintf("http://locahost:%d", port)
}

func respondWith(w http.ResponseWriter, status int, headers Headers) {
  for key, value := range headers {
    w.Header().Set(key, value)
  }

  w.WriteHeader(status)
}

func extractUrl(r *http.Request) string {
  url := make([]byte, r.ContentLength, r.ContentLength)
  r.Body.Read(url)
  return string(url)
}

func Redirecter(w http.ResponseWriter, r *http.Request) {
  path := strings.Split(r.URL.Path, "/")
  id := path[len(path) - 1]

  if url := url.Find(id); url != nil {
    http.Redirect(w, r, url.UrlOriginal, http.StatusMovedPermanently)
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
  respondWith(w, status, Headers{"Location": urlShort})
}

func Home(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello")
}

func main() {
  url.ConfigRepository(url.InitializeRepository())

  http.HandleFunc("/api/shorter", Shorter)
  http.HandleFunc("/r/", Redirecter)
  http.HandleFunc("/home", Home)

  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}