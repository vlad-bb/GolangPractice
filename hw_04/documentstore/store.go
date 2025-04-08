package documentstore

type Store struct {
	storage map[string]Collection
}

func NewStore() *Store {
	return &Store{
		storage: make(map[string]Collection),
	}

}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
	// Якщо ж колекція вже створення, то повертаємо `false` та nil
	_, ok := s.storage[name]
	if ok {
		return false, nil
	}
	collection := Collection{
		cfg:       *cfg,
		documents: make(map[string]Document),
	}
	s.storage[name] = collection

	return true, &collection

}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	collection, ok := s.storage[name]
	if ok {
		return &collection, true
	}

	return nil, false
}

func (s *Store) DeleteCollection(name string) bool {
	_, ok := s.storage[name]
	if ok {
		delete(s.storage, name)

		return true
	}

	return false
}
