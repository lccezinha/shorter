package url

import (
	"gopkg.in/mgo.v2"
	"labix.org/v2/mgo/bson"
	"time"
)

var urlConnection = "localhost:27017"

const database = "shorter"
const collection = "urls"

type mongoRepository struct {
	session *mgo.Session
}

func (mr *mongoRepository) getCollection() *mgo.Collection {
	return mr.session.DB(database).C(collection)
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
	url := &Url{}

	if url = mr.FindById(id); url != nil {
		return true
	}

	return false
}

func (mr *mongoRepository) FindById(id string) *Url {
	url := Url{}

	err := mr.getCollection().Find(bson.M{"id": id}).One(&url)

	if err == nil {
		return &url
	}

	return nil
}

func (mr *mongoRepository) FindByUrl(url string) *Url {
	u := Url{}

	err := mr.getCollection().Find(bson.M{"url_original": url}).One(&u)

	if err == nil {
		return &u
	}

	return nil
}

func (mr *mongoRepository) Save(url Url) error {
	u := &Url{Id: url.Id, CreatedAt: time.Now(), UrlOriginal: url.UrlOriginal}

	err := mr.getCollection().Insert(u)

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
