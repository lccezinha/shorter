package url

type repositoryMemory struct {
  urls map[string]*Url
}

func Initialize() *repositoryMemory {
  return &repositoryMemory{make(map[string]*Url)}
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