package document_store

import (
	"os"
	"testing"
)

func TestStore_CreateAndGetCollection(t *testing.T) {
	store := NewStore()

	ok, col := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if !ok || col == nil {
		t.Fatal("expected collection to be created")
	}

	col2, found := store.GetCollection("users")
	if !found || col2 == nil {
		t.Fatal("expected to find collection")
	}
}

func TestStore_CreateCollection_AlreadyExists(t *testing.T) {
	store := NewStore()
	store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})

	ok, col := store.CreateCollection("users", &CollectionConfig{PrimaryKey: "id"})
	if ok || col != nil {
		t.Fatal("expected collection not to be created again")
	}
}

func TestStore_DeleteCollection(t *testing.T) {
	store := NewStore()
	store.CreateCollection("test", &CollectionConfig{PrimaryKey: "id"})

	ok := store.DeleteCollection("test")
	if !ok {
		t.Fatal("expected deletion to succeed")
	}

	_, found := store.GetCollection("test")
	if found {
		t.Fatal("expected collection to be deleted")
	}

	if store.DeleteCollection("nonexistent") {
		t.Fatal("expected deletion of nonexistent collection to fail")
	}
}

func TestStore_DumpAndLoad(t *testing.T) {
	store := NewStore()
	_, col := store.CreateCollection("docs", &CollectionConfig{PrimaryKey: "id"})
	col.Put(Document{Fields: map[string]DocumentField{
		"id":   {Type: DocumentFieldTypeString, Value: "1"},
		"name": {Type: DocumentFieldTypeString, Value: "Alice"},
	}})

	dump, err := store.Dump()
	if err != nil {
		t.Fatalf("expected dump to succeed, got error: %v", err)
	}

	newStore, err := NewStoreFromDump(dump)
	if err != nil {
		t.Fatalf("expected unmarshal from dump to succeed, got error: %v", err)
	}

	loadedCol, ok := newStore.GetCollection("docs")
	if !ok {
		t.Fatal("expected to find collection after unmarshal")
	}

	doc, found := loadedCol.Get("1")
	if !found || doc.Fields["name"].Value != "Alice" {
		t.Fatal("expected to find document with correct value")
	}
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
	if err != nil {
		t.Fatalf("expected dump to file to succeed, got: %v", err)
	}

	loadedStore, err := NewStoreFromFile(filename)
	if err != nil {
		t.Fatalf("expected load from file to succeed, got: %v", err)
	}

	loadedCol, ok := loadedStore.GetCollection("logs")
	if !ok {
		t.Fatal("expected to find loaded collection")
	}

	doc, found := loadedCol.Get("log1")
	if !found || doc.Fields["msg"].Value != "Hello" {
		t.Fatalf("expected document with correct msg field, got: %v", doc.Fields["msg"].Value)
	}
}
