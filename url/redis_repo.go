package url

import (
  "gopkg.in/redis.v3"
  "strconv"
  "strings"
)

type redisRepository struct {
  client *redis.Client
  urls map[string]*Url
}

func InitializeRedisRepository() *redisRepository {
  client := redis.NewClient(&redis.Options{
    Addr: ":6379",
    Password: "",
    DB: 0,
  })

  return &redisRepository{
    client,
    make(map[string]*Url),
  }
}

func (r *redisRepository) Persisted(id string) bool {
  // key := strings.Join([]string{"urls", url.Id}, ":")

  // r.client.HExists(key)
}

func (r *redisRepository) FindById(id string) *Url {
}

func (r *redisRepository) FindByUrl(url string) *Url {
}

func (r *redisRepository) Save(url Url) error {
  // key := strings.Join([]string{"urls", url.Id}, ":")
  // r.client.HSet(key, url.UrlOriginal)

  // return nil
}

func (r *redisRepository) RegisterClick(id string) {
  r.client.HSet("clicks", id, "1")
}

func (r *redisRepository) ShowClicks(id string) int {
  clicks := r.client.HGet("clicks", id).Val()
  i, err := strconv.Atoi(clicks)

  if err != nil {
    return i
  }

  return 0
}