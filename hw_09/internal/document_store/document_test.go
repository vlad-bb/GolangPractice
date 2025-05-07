package document_store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicDocumentFields(t *testing.T) {
	doc := NewDocumentBuilder().
		WithField("name", DocumentFieldTypeString, "Alice").
		WithField("age", DocumentFieldTypeNumber, 28).
		WithField("active", DocumentFieldTypeBool, true).
		Build()

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
		assert.True(t, ok, "field %s should exist", tt.fieldName)
		assert.Equal(t, tt.wantType, field.Type, "field %s has wrong type", tt.fieldName)
		assert.Equal(t, tt.wantValue, field.Value, "field %s has wrong value", tt.fieldName)
	}
}

func TestComplexTypes(t *testing.T) {
	doc := NewDocumentBuilder().
		WithField("tags", DocumentFieldTypeArray, []string{"go", "backend"}).
		WithField("profile", DocumentFieldTypeObject, map[string]interface{}{"city": "Lviv"}).
		Build()

	tagsField, ok := doc.Fields["tags"]
	assert.True(t, ok, "tags field should exist")
	assert.Equal(t, DocumentFieldTypeArray, tagsField.Type)
	assert.Equal(t, []string{"go", "backend"}, tagsField.Value)

	profileField, ok := doc.Fields["profile"]
	assert.True(t, ok, "profile field should exist")
	assert.Equal(t, DocumentFieldTypeObject, profileField.Type)
	assert.Equal(t, map[string]interface{}{"city": "Lviv"}, profileField.Value)
}
