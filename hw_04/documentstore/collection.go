package documentstore

import "fmt"

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s *Collection) Put(doc Document) {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	pk := s.cfg.PrimaryKey
	pkField, err := doc.Fields[pk]
	if err != true {
		fmt.Printf("Document does not have a field named %s.\n", pk)
		return
	}

	if pkField.Type != DocumentFieldTypeString {
		fmt.Printf("Invalid  key %v, should be 'string'\n", pkField.Type)
		return
	}
	key, ok := pkField.Value.(string)
	if !ok {
		fmt.Printf("Invalid key %v, should be 'string'\n", pkField.Value)
	}
	if key == "" {
		fmt.Printf("Invalid key %v\n", pkField.Value)
	}
	s.documents[key] = doc
}

func (s *Collection) Get(key string) (*Document, bool) {
	doc, ok := s.documents[key]
	if !ok {
		return nil, false
	}
	return &doc, true
}

func (s *Collection) Delete(key string) bool {
	_, ok := s.documents[key]
	if ok {
		delete(s.documents, key)
		return true
	} else {
		return false
	}
}

func (s *Collection) List() []Document {
	docs := make([]Document, 0, len(s.documents))
	for _, doc := range s.documents {
		docs = append(docs, doc)
	}
	return docs
}
