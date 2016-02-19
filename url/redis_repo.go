package url

import (
  "github.com/garyburd/redigo/redis"
  "strconv"
  "log"
  "encoding/json"
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
  persisted, _ := redis.Bool(r.client.Do("hexists", "urls", id))

  return persisted
}

func (r *redisRepository) FindById(id string) *Url {
  var u Url
  data, _ := r.client.Do("hget", "urls", id)

  if data != "" {
    urlJSON, _ := []byte{`data`}

    if err := json.Unmarshal(urlJSON, &u); err != nil {
      log.Fatal(err)
    } else {
      return &u
    }
  }

  return nil
}

func (r *redisRepository) FindByUrl(url string) *Url {
  var u Url
  data, _ := r.client.Do("hget", "urls", url)

  if data != "" {
    urlJSON, _ := []byte{`data`}

    if err := json.Unmarshal(urlJSON, &u); err != nil {
      log.Fatal(err)
    } else {
      return &u
    }
  }

  return nil
}

func (r *redisRepository) Save(url Url) error {
  urlJSON, err := json.Marshal(url)

  if err != nil {
    log.Fatal(err)
  }

  r.client.Do("hset", "urls", url.Id, urlJSON)
  r.client.Do("hset", "urls", url.UrlOriginal, urlJSON)

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