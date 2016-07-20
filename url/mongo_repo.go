package url

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
)

var urlConnection = "localhost:27017"

const database = "shorter"
const collection = "urls"

type mongoRepository struct {
	session *mgo.Session
}

func InitializeMongoRepository() *mongoRepository {
	session, err := mgo.Dial(urlConnection)

	if err != nil {
		panic("Database connection Error!")
	}

	session.SetMode(mgo.Monotonic, true)

	return &mongoRepository{session}
}

func (mr *mongoRepository) Persisted(id string) bool {
	return false
	// var u Url

	// mr.connection.DB(database).C(collection)
}

func (mr *mongoRepository) FindById(id string) *Url {
	return nil
}

func (mr *mongoRepository) FindByUrl(url string) *Url {
	return nil
}

func (mr *mongoRepository) Save(url Url) error {
	u := &Url{id: url.Id, created_at: time.Now(), url_original: url.UrlOriginal}

	err := mr.session.DB(database).C(collection).Insert(u)

	if err != nil {
		panic(err)
	}

	return nil
}

func (mr *mongoRepository) RegisterClick(id string) {
}

func (mr *mongoRepository) ShowClicks(id string) int {
	return 1
}
