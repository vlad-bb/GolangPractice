package document_store

import (
	"encoding/json"
	"testing"
)

func TestCollection_PutAndGet(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})

	doc := Document{
		Fields: map[string]DocumentField{
			"id": {
				Type:  DocumentFieldTypeString,
				Value: "123",
			},
			"name": {
				Type:  DocumentFieldTypeString,
				Value: "Test",
			},
		},
	}

	err := c.Put(doc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, ok := c.Get("123")
	if !ok {
		t.Fatalf("expected document to exist")
	}
	if got.Fields["name"].Value != "Test" {
		t.Errorf("expected name to be Test, got %v", got.Fields["name"].Value)
	}
}

func TestCollection_Put_Errors(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})

	// Missing primary key
	doc1 := Document{Fields: map[string]DocumentField{}}
	err := c.Put(doc1)
	if err == nil {
		t.Error("expected error for missing primary key")
	}

	// Wrong type
	doc2 := Document{Fields: map[string]DocumentField{
		"id": {Type: DocumentFieldTypeNumber, Value: 123},
	}}
	err = c.Put(doc2)
	if err == nil {
		t.Error("expected error for wrong primary key type")
	}

	// Invalid value (not string)
	doc3 := Document{Fields: map[string]DocumentField{
		"id": {Type: DocumentFieldTypeString, Value: 123},
	}}
	err = c.Put(doc3)
	if err == nil {
		t.Error("expected error for invalid primary key value")
	}
}

func TestCollection_Delete(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})
	doc := Document{
		Fields: map[string]DocumentField{
			"id": {Type: DocumentFieldTypeString, Value: "abc"},
		},
	}
	_ = c.Put(doc)

	ok := c.Delete("abc")
	if !ok {
		t.Error("expected document to be deleted")
	}

	_, exists := c.Get("abc")
	if exists {
		t.Error("expected document to be gone after deletion")
	}

	// Deleting non-existent key
	if c.Delete("not-there") {
		t.Error("expected false when deleting non-existent key")
	}
}

func TestCollection_List(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})

	c.Put(Document{Fields: map[string]DocumentField{
		"id":   {Type: DocumentFieldTypeString, Value: "1"},
		"name": {Type: DocumentFieldTypeString, Value: "First"},
	}})
	c.Put(Document{Fields: map[string]DocumentField{
		"id":   {Type: DocumentFieldTypeString, Value: "2"},
		"name": {Type: DocumentFieldTypeString, Value: "Second"},
	}})

	docs := c.List()
	if len(docs) != 2 {
		t.Errorf("expected 2 documents, got %d", len(docs))
	}
}

func TestCollection_JSON_MarshalUnmarshal(t *testing.T) {
	c := NewCollection(CollectionConfig{PrimaryKey: "id"})
	_ = c.Put(Document{
		Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: "1"},
			"name": {Type: DocumentFieldTypeString, Value: "Doc1"},
		},
	})

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var newC Collection
	if err := json.Unmarshal(data, &newC); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	doc, ok := newC.Get("1")
	if !ok {
		t.Fatal("expected document to be present after unmarshal")
	}
	if doc.Fields["name"].Value != "Doc1" {
		t.Errorf("expected name to be Doc1, got %v", doc.Fields["name"].Value)
	}
}
