package url

import (
  "math/rand"
  "net/url"
  "time"
)

const (
  size = 5
  symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-+"
)

var repository Repository

type Repository interface {
  Persisted(id string) bool
  FindById(id string) *Url
  FindByUrl(url string) *Url
  Save(url Url) error
}

type Url struct {
  Id int
  CreatedAt time.Time
  UrlOriginal string
}

func ConfigRepository(r Repository) {
  repository = r
}

func init() {
  rand.Seed(time.Now().UnixNano())
}

func generateId() string {
  newId := func() string {
    id := make([]byte, size, size)

    for i := range id {
      id[i] = symbols[rand.Intn(len(symbols))]
    }

    return string(id)
  }

  for {
    id := newId(); !repo.Persisted(id) {
      return id
    }
  }
}

func Find(id string) *Url {
  return repo.FindById(id)
}

func FetchUrl(urlOriginal string) (url *string, isNew bool, err error) {
  if u = repo.FindByUrl(urlOriginal); u != nil {
    return u, false, nil
  }

  if _, err = url.ParseRequestURI(urlOriginal); err != nil {
    return nil, false, err
  }

  url := Url{generateId(), time.Now(), urlOriginal}
  repo.Save(url)

  return &url, true, nil
}