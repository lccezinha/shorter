package url

import (
  "gopkg.in/redis.v3"
  "strconv"
  "encoding/json"
  "fmt"
)

type redisRepository struct {
  client *redis.Client
}

func InitializeRedisRepository() *redisRepository {
  client := redis.NewClient(&redis.Options{
    Addr: ":6379",
    Password: "",
    DB: 0,
  })

  return &redisRepository{client}
}

func (r *redisRepository) Persisted(id string) bool {
  persisted := r.client.HExists("urls", id)

  if persisted != nil {
    return true
  } else {
    return false
  }
}

func (r *redisRepository) FindById(id string) *Url {
  // fmt.Println(id)
  // var u Url
  // data := r.client.HMGet("urls", id).Val()
  // fmt.Println(data)

  // if data != nil {
  //   urlJSON, _ := []byte{string(data)}
  //   err := json.Unmarshal(urlJSON, &u)

  //   if err != nil {
  //     return &u
  //   }
  // }

  return &Url{Id: "123qw", UrlOriginal: "http://www.pudim.com.br"}
}

func (r *redisRepository) FindByUrl(url string) *Url {
  // var u Url
  // data := r.client.HMGet("urls", url).Val()

  // if data != nil {
  //   urlJSON, _ := []byte{string(data)}
  //   err := json.Unmarshal(urlJSON, &u)

  //   if err != nil {
  //     return &u
  //   }
  // }

  return &Url{Id: "123qw", UrlOriginal: "http://www.pudim.com.br"}
}

func (r *redisRepository) Save(url Url) error {
  fmt.Println(url)
  urlJSON, err := json.Marshal(url)

  fmt.Println(string(urlJSON))

  if err != nil {
    return err
  }

  r.client.HMSet("urls", url.Id, string(urlJSON))
  r.client.HMSet("urls", url.UrlOriginal, string(urlJSON))

  return nil
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