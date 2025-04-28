package document_store

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestStore_CreateAndGetCollection(t *testing.T) {
	store := NewStore()

	ok, col := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	assert.True(t, ok, "expected collection to be created")
	assert.NotNil(t, col, "expected created collection not to be nil")

	col2, found := store.GetCollection("users")
	assert.True(t, found, "expected to find collection")
	assert.NotNil(t, col2, "expected found collection not to be nil")
}

func TestStore_CreateCollection_AlreadyExists(t *testing.T) {
	store := NewStore()
	_, _ = store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

	ok, col := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	assert.False(t, ok, "expected collection not to be created again")
	assert.Nil(t, col, "expected no collection to be returned on duplicate create")
}

func TestStore_DeleteCollection(t *testing.T) {
	store := NewStore()
	store.CreateCollection("test", &CollectionConfig{PrimaryKey: "id"})

	assert.True(t, store.DeleteCollection("test"), "expected deletion to succeed")
	_, found := store.GetCollection("test")
	assert.False(t, found, "expected collection to be deleted")
	assert.False(t, store.DeleteCollection("nonexistent"), "expected deletion of nonexistent collection to fail")
}

func TestStore_DumpAndLoad(t *testing.T) {
	store := NewStore()
	_, col := store.CreateCollection("docs", &CollectionConfig{PrimaryKey: "id"})
	col.Put(Document{Fields: map[string]DocumentField{
		"id":   {Type: DocumentFieldTypeString, Value: "1"},
		"name": {Type: DocumentFieldTypeString, Value: "Alice"},
	}})

	dump, err := store.Dump()
	assert.NoError(t, err, "expected dump to succeed")

	newStore, err := NewStoreFromDump(dump)
	assert.NoError(t, err, "expected unmarshal from dump to succeed")

	loadedCol, ok := newStore.GetCollection("docs")
	assert.True(t, ok, "expected to find collection after unmarshal")

	doc, found := loadedCol.Get("1")
	assert.True(t, found, "expected to find document with correct value")
	assert.Equal(t, "Alice", doc.Fields["name"].Value, "expected document to have correct name")
}

func TestStore_DumpToFile_And_NewStoreFromFile(t *testing.T) {
	filename := "test_store_dump.gz"
	defer os.Remove(filename)

	store := NewStore()
	_, col := store.CreateCollection("logs", &CollectionConfig{PrimaryKey: "id"})
	col.Put(Document{Fields: map[string]DocumentField{
		"id":  {Type: DocumentFieldTypeString, Value: "log1"},
		"msg": {Type: DocumentFieldTypeString, Value: "Hello"},
	}})
	
	err := store.DumpToFile(filename)
	assert.NoError(t, err, "expected dump to file to succeed")

	loadedStore, err := NewStoreFromFile(filename)
	assert.NoError(t, err, "expected load from file to succeed")

	loadedCol, ok := loadedStore.GetCollection("logs")
	assert.True(t, ok, "expected to find loaded collection")

	doc, found := loadedCol.Get("log1")
	assert.True(t, found, "expected to find document with correct msg field")
	assert.Equal(t, "Hello", doc.Fields["msg"].Value, "expected document to have correct msg field")
}
