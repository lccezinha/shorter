package main

import (
  "fmt"
  "log"
  "net/http"
  "strings"

  "github.com/lccezinha/shorter/url"
)

type Headers map[string]string

/*
Ainda é necessário implementar os exemplos do começo do capítulo (métodos do pacote URL)
*/

var (
  port int
  urlBase string
)

func init() {
  port = "8080"
  urlBase = fmt.Printf("http://locahost:%d", urlBase)
}

func respondWith(w http.ResponseWriter, status int, headers Headers) {
  for key, value := range headers {
    w.Header.Set(key, value)
  }

  w.WriteHeader(status)
}

func extracUrl(r *http.Request) string {
  url := make([]byte, r.ContentLength, r.ContentLength)
  r.Body.Read(url)
  return string(url)
}

func Shorter(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    respondWith(w, http.StatusMethodNotAllowed, Headers{ "Allow":"POST", })
    return
  }
}

func main() {
  http.HandleFunc("/api/short", Shorter)
  http.HandleFunc("/r/", Redirecter)

  log.Fatal(
    http.ListenAndServe(fmt.Printf(":%d", port), nil)
  )
}