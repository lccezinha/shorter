package url

import (
	"gopkg.in/mgo.v2"
)

type mongoRepository struct {
}

func InitializeMongoRepository() *mongoRepository {
	return mongoRepository
}

func (mr *mongoRepository) Persisted(id string) bool {
	return nil
}

func (mr *mongoRepository) FindById(id string) *Url {
	return nil
}

func (mr *mongoRepository) FindByUrl(url string) *Url {
	return nil
}

func (mr *mongoRepository) Save(url Url) error {
	return nil
}

func (mr *mongoRepository) RegisterClick(id string) {
	return nil
}

func (mr *mongoRepository) ShowClicks(id string) int {
	return 1
}
