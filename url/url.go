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
  RegisterClick(id string)
  ShowClicks(id string) int
}

type Url struct {
  Id          string    `json:"id"`
  CreatedAt   time.Time `json:"created_at"`
  UrlOriginal string    `json:"url_original"`
}

type Stats struct {
  Url   *Url `json:"url"`
  Clicks int `json:"clicks"`
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
    if id := newId(); !repository.Persisted(id) {
      return id
    }
  }
}

func Find(id string) *Url {
  return repository.FindById(id)
}

func RegisterClick(id string) {
  repository.RegisterClick(id)
}

func FetchUrl(urlOriginal string) (u *Url, isNew bool, err error) {
  if u = repository.FindByUrl(urlOriginal); u != nil {
    return u, false, nil
  }

  if _, err = url.ParseRequestURI(urlOriginal); err != nil {
    return nil, false, err
  }

  url := Url{generateId(), time.Now(), urlOriginal}
  repository.Save(url)

  return &url, true, nil
}

func (u *Url) ShowStats() *Stats {
  clicks := repository.ShowClicks(u.Id)

  return &Stats{u, clicks}
}