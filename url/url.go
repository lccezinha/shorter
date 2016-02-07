package url

import (
  "math/rand"
  "net/url"
  "time"
)

func init() {
  rand.Seed(time.Now().UnixNano())
}

type Url struct {
  Id int
  CreatedAt time.Time
  urlOriginal string
}

const (
  size = 5
  symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-+"
)