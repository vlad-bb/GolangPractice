package document_store

import "fmt"

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

func (s *Collection) Put(doc Document) error {
	pk := s.cfg.PrimaryKey

	pkField, ok := doc.Fields[pk]
	if !ok {
		return fmt.Errorf("%w: %s", ErrPrimaryKeyNotFound, pk)
	}

	if pkField.Type != DocumentFieldTypeString {
		return fmt.Errorf("%w: %v", ErrPrimaryKeyWrongType, pkField.Type)
	}

	key, ok := pkField.Value.(string)
	if !ok || key == "" {
		return fmt.Errorf("%w: %v", ErrPrimaryKeyInvalidValue, pkField.Value)
	}

	s.documents[key] = doc
	return nil
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
