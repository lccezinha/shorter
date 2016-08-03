package url

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

var mongoConfig MongoConfig

func init() {
	file, err := ioutil.ReadFile("./mongo_config.json")

	if err != nil {
		fmt.Println(file)
		log.Fatal(err)
	}

	fmt.Println(file)

	mongoConfig = MongoConfig{}
	json.Unmarshal(file, &mongoConfig)
}

type mongoRepository struct {
	session *mgo.Session
}

func (mr *mongoRepository) getCollection() *mgo.Collection {
	return mr.session.DB(mongoConfig.Database).C(mongoConfig.Collection)
}

func InitializeMongoRepository() *mongoRepository {
	session, err := mgo.Dial(mongoConfig.UrlConnection)

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
	u := &Url{Id: url.Id, CreatedAt: time.Now(), UrlOriginal: url.UrlOriginal, Clicks: 0}

	err := mr.getCollection().Insert(u)

	if err != nil {
		panic(err)
	}

	return nil
}

func (mr *mongoRepository) RegisterClick(id string) {
	url := &Url{}

	if url = mr.FindById(id); url != nil {
		change := mgo.Change{
			Update: bson.M{"$inc": bson.M{"clicks": 1}},
		}

		mr.getCollection().Find(bson.M{"id": id}).Apply(change, &url)
	}
}

func (mr *mongoRepository) ShowStats(id string) *Url {
	url := mr.FindById(id)

	return url
}
