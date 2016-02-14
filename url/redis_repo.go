package url

type redisRepository struct{}

func InitializeRedisRepository() *redisRepository {
  return &redisRepository{}
}

func (r *redisRepository) Persisted(id string) bool {
}

func (r *redisRepository) FindById(id string) *Url {
}

func (r *redisRepository) FindByUrl(url string) *Url {
}

func (r *redisRepository) Save(url Url) error {
}

func (r *redisRepository) RegisterClick(id string) {
}

func (r *redisRepository) ShowClicks(id string) int {
}