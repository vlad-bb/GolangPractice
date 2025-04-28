package document_store

import (
	"encoding/json"
	"fmt"
)

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
}

type CollectionConfig struct {
	PrimaryKey string
}

func NewCollection(cfg CollectionConfig) *Collection {
	return &Collection{
		cfg:       cfg,
		documents: make(map[string]Document),
	}
}

func (c *Collection) Put(doc Document) error {
	pk := c.cfg.PrimaryKey

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

	c.documents[key] = doc
	return nil
}

func (c *Collection) Get(key string) (*Document, bool) {
	doc, ok := c.documents[key]
	if !ok {
		return nil, false
	}

	return &doc, true
}

func (c *Collection) Delete(key string) bool {
	_, ok := c.documents[key]
	if ok {
		delete(c.documents, key)

		return true
	} else {
		return false
	}
}

func (c *Collection) List() []Document {
	docs := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		docs = append(docs, doc)
	}

	return docs
}

func (c *Collection) MarshalJSON() ([]byte, error) {
	cBytes, err := json.Marshal(&struct {
		Cfg       CollectionConfig    `json:"cfg"`
		Documents map[string]Document `json:"documents"`
	}{
		Cfg:       c.cfg,
		Documents: c.documents,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMarshalJSONFailed, err)
	}
	return cBytes, nil
}

func (c *Collection) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Cfg       CollectionConfig    `json:"cfg"`
		Documents map[string]Document `json:"documents"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshalJSONFailed, err)
	}

	c.cfg = aux.Cfg
	c.documents = aux.Documents
	return nil
}
