package url

type repositoryMemory struct {
  urls map[string]*Url
  clicks map[string]int
}

func InitializeRepository() *repositoryMemory {
  return &repositoryMemory{
    make(map[string]*Url),
    make(map[string]int),
  }
}

func (r *repositoryMemory) Persisted(id string) bool {
  _, persisted := r.urls[id]

  return persisted
}

func (r *repositoryMemory) FindById(id string) *Url {
  return r.urls[id]
}

func (r *repositoryMemory) FindByUrl(url string) *Url {
  for _, u := range r.urls {
    if u.UrlOriginal == url {
      return u
    }
  }

  return nil
}

func (r *repositoryMemory) Save(url Url) error {
  r.urls[url.Id] = &url
  return nil
}

func (r *repositoryMemory) RegisterClick(id string) {
  r.clicks[id] += 1
}

func (r *repositoryMemory) ShowClicks(id string) int {
  return r.clicks[id]
}