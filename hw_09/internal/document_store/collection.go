package document_store

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Index struct {
	fieldName string
	values    map[string][]string
}

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
	indexes   map[string]*Index
}

type CollectionConfig struct {
	PrimaryKey string
}

type QueryParams struct {
	Desc     bool
	MinValue *string
	MaxValue *string
}

func NewCollection(cfg CollectionConfig) *Collection {
	return &Collection{
		cfg:       cfg,
		documents: make(map[string]Document),
		indexes:   make(map[string]*Index),
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

	if oldDoc, exists := c.documents[key]; exists {
		c.removeDocumentFromIndexes(key, oldDoc)
	}
	c.documents[key] = doc
	c.addDocumentToIndexes(key, doc)
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
	doc, ok := c.documents[key]
	if ok {
		c.removeDocumentFromIndexes(key, doc)
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

type serializedIndex struct {
	FieldName string              `json:"fieldName"`
	Values    map[string][]string `json:"values"`
}

func (idx *Index) MarshalJSON() ([]byte, error) {
	return json.Marshal(&serializedIndex{
		FieldName: idx.fieldName,
		Values:    idx.values,
	})
}

func (idx *Index) UnmarshalJSON(data []byte) error {
	aux := &serializedIndex{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	idx.fieldName = aux.FieldName
	idx.values = aux.Values
	return nil
}

func (c *Collection) MarshalJSON() ([]byte, error) {
	cBytes, err := json.Marshal(&struct {
		Cfg       CollectionConfig    `json:"cfg"`
		Documents map[string]Document `json:"documents"`
		Indexes   map[string]*Index   `json:"indexes"`
	}{
		Cfg:       c.cfg,
		Documents: c.documents,
		Indexes:   c.indexes,
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
		Indexes   map[string]*Index   `json:"indexes"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshalJSONFailed, err)
	}

	c.cfg = aux.Cfg
	c.documents = aux.Documents
	c.indexes = aux.Indexes
	if c.indexes == nil {
		c.indexes = make(map[string]*Index)
	}
	return nil
}

func (c *Collection) CreateIndex(fieldName string) error {
	if _, exists := c.indexes[fieldName]; exists {
		return ErrIndexAlreadyExists
	}
	index := &Index{
		fieldName: fieldName,
		values:    make(map[string][]string),
	}

	for key, doc := range c.documents {
		field, ok := doc.Fields[fieldName]
		if !ok {
			continue
		}

		if field.Type != DocumentFieldTypeString {
			continue
		}

		strValue, ok := field.Value.(string)
		if !ok {
			continue
		}

		index.values[strValue] = append(index.values[strValue], key)
	}

	c.indexes[fieldName] = index
	return nil
}

func (c *Collection) DeleteIndex(fieldName string) error {
	if _, exists := c.indexes[fieldName]; !exists {
		return ErrIndexNotFound
	}

	delete(c.indexes, fieldName)
	return nil
}

func (c *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	index, exists := c.indexes[fieldName]
	if !exists {
		return nil, ErrIndexNotFound
	}

	values := make([]string, 0, len(index.values))
	for value := range index.values {
		values = append(values, value)
	}

	filteredValues := filterValues(values, params.MinValue, params.MaxValue)

	if params.Desc {
		sort.Sort(sort.Reverse(sort.StringSlice(filteredValues)))
	} else {
		sort.Strings(filteredValues)
	}

	result := make([]Document, 0)
	for _, value := range filteredValues {
		docKeys := index.values[value]

		for _, key := range docKeys {
			doc, exists := c.documents[key]
			if exists {
				result = append(result, doc)
			}
		}
	}

	return result, nil
}

func (c *Collection) addDocumentToIndexes(key string, doc Document) {
	for fieldName, index := range c.indexes {
		field, ok := doc.Fields[fieldName]
		if !ok {
			continue
		}

		if field.Type != DocumentFieldTypeString {
			continue
		}

		strValue, ok := field.Value.(string)
		if !ok {
			continue
		}

		index.values[strValue] = append(index.values[strValue], key)
	}
}

func (c *Collection) removeDocumentFromIndexes(key string, doc Document) {
	for fieldName, index := range c.indexes {
		field, ok := doc.Fields[fieldName]
		if !ok {
			continue
		}

		if field.Type != DocumentFieldTypeString {
			continue
		}

		strValue, ok := field.Value.(string)
		if !ok {
			continue
		}

		docKeys := index.values[strValue]
		for i, docKey := range docKeys {
			if docKey == key {
				index.values[strValue] = append(docKeys[:i], docKeys[i+1:]...)
				break
			}
		}

		if len(index.values[strValue]) == 0 {
			delete(index.values, strValue)
		}
	}
}

func filterValues(values []string, minValue, maxValue *string) []string {
	if minValue == nil && maxValue == nil {
		return values
	}

	filtered := make([]string, 0, len(values))
	for _, v := range values {
		if minValue != nil && v < *minValue {
			continue
		}
		if maxValue != nil && v > *maxValue {
			continue
		}
		filtered = append(filtered, v)
	}
	return filtered
}
