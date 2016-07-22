package url

import (
	"labix.org/v2/mgo/bson"
	"math/rand"
	"net/url"
	"time"
)

const (
	size    = 5
	symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-+"
)

var repository Repository

type Repository interface {
	Persisted(id string) bool
	FindById(id string) *Url
	FindByUrl(url string) *Url
	Save(url Url) error
	RegisterClick(id string)
	ShowStats(id string) *Url
}

type Url struct {
	Mid         bson.ObjectId `bson:"_id,omitempty"`
	Id          string        `bson:"id,omitempty" json:"id"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UrlOriginal string        `bson:"url_original" json:"url_original"`
	Clicks      int           `bson:"clicks" json:"clicks"`
}

type Stats struct {
	Url *Url `json:"url"`
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

	url := Url{Id: generateId(), CreatedAt: time.Now(), UrlOriginal: urlOriginal}
	repository.Save(url)

	return &url, true, nil
}

func (u *Url) ShowStats() *Stats {
	url := repository.FindById(u.Id)

	if url != nil {
		return &Stats{url}
	}

	return nil
}
