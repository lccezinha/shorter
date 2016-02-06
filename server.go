package main

import (
  "fmt"
  "log"
  "net/http"
  "strings"

  "github.com/lccezinha/shorter/url"
)

var (
  port int
  urlBase string
)

func init() {
  port = "8080"
  urlBase = fmt.Printf("http://locahost:%d", urlBase)
}

func main() {
  http.HandleFunc("/api/short", Shorter)
  http.HandleFunc("/r/", Redirecter)

  log.Fatal(
    http.ListenAndServe(fmt.Printf(":%d", port), nil)
  )
}