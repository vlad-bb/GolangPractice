package document_store

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTestDocument(id, name string) Document {
	return Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: id},
			"name": {Type: DocumentFieldTypeString, Value: name},
		},
	}
}

func TestCollection_Indexing(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})

	// Додати індекс за полем "name"
	c.indexes["name"] = &Index{
		values: make(map[string]map[string]struct{}),
	}

	// Додати документ
	doc := newTestDocument("42", "John")
	err := c.Put(doc)
	assert.NoError(t, err)

	// Перевірити, чи індекс "name" містить "John" -> "42"
	index, exists := c.indexes["name"]
	assert.True(t, exists, "index 'name' should exist")

	index.mu.RLock()
	defer index.mu.RUnlock()

	docKeys, ok := index.values["John"]
	assert.True(t, ok, "index for 'John' should exist")
	_, keyExists := docKeys["42"]
	assert.True(t, keyExists, "document ID should be indexed under 'John'")
}

func TestCollection_PutAndGet(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})

	doc := newTestDocument("123", "Test")

	err := c.Put(doc)
	assert.NoError(t, err, "expected no error on Put")

	got, ok := c.Get("123")
	assert.True(t, ok, "expected document to exist")

	assert.Equal(t, "Test", got.Fields["name"].Value, "expected name to be 'Test'")
}

func TestCollection_Put_Errors(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})

	// Missing primary key
	doc1 := Document{Fields: map[string]DocumentField{
		"name": {Type: DocumentFieldTypeString, Value: "test"},
	}}
	err := c.Put(doc1)
	assert.Error(t, err, "expected error for missing primary key")

	// Wrong type
	doc2 := Document{Fields: map[string]DocumentField{
		"id": {Type: DocumentFieldTypeNumber, Value: 123},
	}}
	err = c.Put(doc2)
	assert.Error(t, err, "expected error for wrong primary key type")

	// Invalid value (not a string)
	doc3 := Document{Fields: map[string]DocumentField{
		"id": {Type: DocumentFieldTypeString, Value: 123},
	}}
	err = c.Put(doc3)
	assert.Error(t, err, "expected error for invalid primary key value")
}

func TestCollection_Delete(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})
	doc := Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "abc"},
		},
	}
	err := c.Put(doc)
	assert.NoError(t, err, "expected no error on Put")

	ok := c.Delete("abc")
	assert.True(t, ok, "expected document to be deleted")

	_, exists := c.Get("abc")
	assert.False(t, exists, "expected document to be gone after deletion")

	// Deleting non-existent key
	ok = c.Delete("not-there")
	assert.False(t, ok, "expected false when deleting non-existent key")
}

func TestCollection_List(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})

	err := c.Put(newTestDocument("1", "First"))
	assert.NoError(t, err)

	err = c.Put(newTestDocument("2", "Second"))
	assert.NoError(t, err)

	docs := c.List()
	assert.Len(t, docs, 2, "expected 2 documents in the list")
}

func TestCollection_JSON_MarshalUnmarshal(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})
	err := c.Put(newTestDocument("1", "Doc1"))
	assert.NoError(t, err)

	data, err := json.Marshal(c)
	assert.NoError(t, err, "marshal should succeed")

	var newC Collection
	err = json.Unmarshal(data, &newC)
	assert.NoError(t, err, "unmarshal should succeed")

	doc, ok := newC.Get("1")
	assert.True(t, ok, "document should exist after unmarshal")
	assert.Equal(t, "Doc1", doc.Fields["name"].Value)
}
