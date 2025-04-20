package document_store

import (
	"reflect"
	"testing"
)

func TestBasicDocumentFields(t *testing.T) {
	doc := Document{
		Fields: map[string]DocumentField{
			"name": {
				Type:  DocumentFieldTypeString,
				Value: "Alice",
			},
			"age": {
				Type:  DocumentFieldTypeNumber,
				Value: 28,
			},
			"active": {
				Type:  DocumentFieldTypeBool,
				Value: true,
			},
		},
	}

	tests := []struct {
		fieldName string
		wantType  DocumentFieldType
		wantValue interface{}
	}{
		{"name", DocumentFieldTypeString, "Alice"},
		{"age", DocumentFieldTypeNumber, 28},
		{"active", DocumentFieldTypeBool, true},
	}

	for _, tt := range tests {
		field, ok := doc.Fields[tt.fieldName]
		if !ok {
			t.Errorf("field %s not found", tt.fieldName)
			continue
		}
		if field.Type != tt.wantType {
			t.Errorf("field %s has wrong type: got %s, want %s", tt.fieldName, field.Type, tt.wantType)
		}
		if field.Value != tt.wantValue {
			t.Errorf("field %s has wrong value: got %v, want %v", tt.fieldName, field.Value, tt.wantValue)
		}
	}
}

func TestComplexTypes(t *testing.T) {
	doc := Document{
		Fields: map[string]DocumentField{
			"tags": {
				Type:  DocumentFieldTypeArray,
				Value: []string{"go", "backend"},
			},
			"profile": {
				Type: DocumentFieldTypeObject,
				Value: map[string]interface{}{
					"city": "Lviv",
				},
			},
		},
	}

	tags := doc.Fields["tags"]
	wantTags := []string{"go", "backend"}
	if !reflect.DeepEqual(tags.Value, wantTags) {
		t.Errorf("expected tags to be %v, got %v", wantTags, tags.Value)
	}

	profile := doc.Fields["profile"]
	wantProfile := map[string]interface{}{"city": "Lviv"}
	if !reflect.DeepEqual(profile.Value, wantProfile) {
		t.Errorf("expected profile to be %v, got %v", wantProfile, profile.Value)
	}
}
