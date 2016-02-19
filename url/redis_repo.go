package url

import (
  "github.com/garyburd/redigo/redis"
  "strconv"
  "log"
)

type redisRepository struct {
  urls map[string]*Url
  client redis.Conn
}

func InitializeRedisRepository() *redisRepository {
  client, err := redis.Dial("tcp", ":6379")

  if err != nil {
    log.Fatal(err)
  }

  return &redisRepository{make(map[string]*Url), client}
}

func (r *redisRepository) Persisted(id string) bool {
  _, persisted := r.urls[id]

  return persisted
}

func (r *redisRepository) FindById(id string) *Url {
  return r.urls[id]
}

func (r *redisRepository) FindByUrl(url string) *Url {
  for _, u := range r.urls {
    if u.UrlOriginal == url {
      return u
    }
  }

  return nil
}

func (r *redisRepository) Save(url Url) error {
  r.urls[url.Id] = &url
  return nil
}

func (r *redisRepository) RegisterClick(id string) {
  r.client.Do("hincrby", "clicks", id, 1)
}

func (r *redisRepository) ShowClicks(id string) int {
  clicks, err := redis.String(r.client.Do("hget", "clicks", id))

  if err != nil {
    log.Fatal(err)
  }

  i, err := strconv.Atoi(clicks)

  if err != nil {
    log.Fatal(err)
  }

  return i
}