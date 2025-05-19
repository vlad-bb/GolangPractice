package document_store

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

type Index struct {
	fieldName  string
	values     map[string]map[string]struct{}
	sortedKeys []string
	mu         sync.RWMutex
}

type Collection struct {
	cfg       CollectionConfig
	documents map[string]Document
	indexes   map[string]*Index
	mu        sync.RWMutex
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
		mu:        sync.RWMutex{},
	}
}

func (c *Collection) Put(doc Document) error {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	c.mu.RLock()
	defer c.mu.RUnlock()
	doc, ok := c.documents[key]
	if !ok {
		return nil, false
	}

	return &doc, true
}

func (c *Collection) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	c.mu.RLock()
	defer c.mu.RUnlock()
	docs := make([]Document, 0, len(c.documents))
	for _, doc := range c.documents {
		docs = append(docs, doc)
	}

	return docs
}

type serializedIndex struct {
	FieldName  string                         `json:"fieldName"`
	Values     map[string]map[string]struct{} `json:"values"`
	SortedKeys []string                       `json:"sortedKeys"`
}

func (idx *Index) MarshalJSON() ([]byte, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	return json.Marshal(&serializedIndex{
		FieldName:  idx.fieldName,
		Values:     idx.values,
		SortedKeys: idx.sortedKeys,
	})
}

func (idx *Index) UnmarshalJSON(data []byte) error {
	aux := &serializedIndex{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.fieldName = aux.FieldName
	idx.values = aux.Values
	idx.sortedKeys = aux.SortedKeys
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
	c.mu.Lock()
	if _, exists := c.indexes[fieldName]; exists {
		c.mu.Unlock()
		return ErrIndexAlreadyExists
	}
	index := &Index{
		fieldName:  fieldName,
		values:     make(map[string]map[string]struct{}),
		sortedKeys: make([]string, 0),
	}

	documentsCopy := make(map[string]Document, len(c.documents))
	for k, v := range c.documents {
		documentsCopy[k] = v
	}
	c.mu.Unlock()

	for key, doc := range documentsCopy {
		field, ok := doc.Fields[fieldName]
		if !ok || field.Type != DocumentFieldTypeString {
			continue
		}

		strValue, ok := field.Value.(string)
		if !ok {
			continue
		}
		index.mu.Lock()
		if index.values[strValue] == nil {
			index.values[strValue] = make(map[string]struct{})
			index.sortedKeys = insertSorted(index.sortedKeys, strValue)
		}
		index.values[strValue][key] = struct{}{}
		index.mu.Unlock()
	}

	c.mu.Lock()
	c.indexes[fieldName] = index
	c.mu.Unlock()
	return nil
}

func (c *Collection) DeleteIndex(fieldName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.indexes[fieldName]; !exists {
		return ErrIndexNotFound
	}

	delete(c.indexes, fieldName)
	return nil
}

func (c *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	c.mu.RLock()
	index, exists := c.indexes[fieldName]
	c.mu.RUnlock()
	if !exists {
		return nil, ErrIndexNotFound
	}
	index.mu.RLock()
	values := filterSortedValues(index.sortedKeys, params.MinValue, params.MaxValue, params.Desc)
	index.mu.RUnlock()
	result := make([]Document, 0)
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, value := range values {
		docKeys := index.values[value]
		for key := range docKeys {
			if doc, exists := c.documents[key]; exists {
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
		index.mu.Lock()
		if _, exists := index.values[strValue]; !exists {
			index.values[strValue] = make(map[string]struct{})
		}

		index.values[strValue][key] = struct{}{}
		index.mu.Unlock()
	}
}

func (c *Collection) removeDocumentFromIndexes(key string, doc Document) {

	for fieldName, index := range c.indexes {
		field, ok := doc.Fields[fieldName]
		if !ok || field.Type != DocumentFieldTypeString {
			continue
		}

		strValue, ok := field.Value.(string)
		if !ok {
			continue
		}

		index.mu.Lock()
		if docKeys, exists := index.values[strValue]; exists {
			delete(docKeys, key)
			if len(docKeys) == 0 {
				delete(index.values, strValue)
			}
		}
		index.mu.Unlock()
	}
}

func filterSortedValues(values []string, minValue, maxValue *string, desc bool) []string {
	start := 0
	end := len(values)

	if minValue != nil {
		start = sort.Search(len(values), func(i int) bool {
			return values[i] >= *minValue
		})
	}

	if maxValue != nil {
		end = sort.Search(len(values), func(i int) bool {
			return values[i] > *maxValue
		})
	}

	sliced := values[start:end]
	if desc {
		for i, j := 0, len(sliced)-1; i < j; i, j = i+1, j-1 {
			sliced[i], sliced[j] = sliced[j], sliced[i]
		}
	}
	return sliced
}

func insertSorted(slice []string, value string) []string {
	i := sort.SearchStrings(slice, value)
	if i < len(slice) && slice[i] == value {
		return slice // вже існує
	}
	slice = append(slice, "")
	copy(slice[i+1:], slice[i:])
	slice[i] = value
	return slice
}
