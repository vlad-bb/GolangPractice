package document_store

import (
	"encoding/json"
	"fmt"
)

func MarshalDocument(input any) (*Document, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMarshalJSONFailed, err)
	}

	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalJSONFailed, err)
	}

	doc := &Document{
		Fields: make(map[string]DocumentField),
	}

	for key, val := range rawMap {
		var fieldType DocumentFieldType
		switch v := val.(type) {
		case string:
			fieldType = DocumentFieldTypeString
		case float64:
			fieldType = DocumentFieldTypeNumber
		case bool:
			fieldType = DocumentFieldTypeBool
		case []interface{}:
			fieldType = DocumentFieldTypeArray
		case map[string]interface{}:
			fieldType = DocumentFieldTypeObject
		default:
			return nil, fmt.Errorf("%w: key=%s, type=%T", ErrUnsupportedFieldType, key, v)
		}
		doc.Fields[key] = DocumentField{
			Type:  fieldType,
			Value: val,
		}
	}

	return doc, nil
}

func UnmarshalDocument(doc *Document, output any) error {
	if doc == nil {
		return ErrDocumentNil
	}
	if doc.Fields == nil {
		return fmt.Errorf("document fields are nil: %w", ErrUnsupportedDocumentField)
	}

	tmpMap := make(map[string]interface{})
	for key, field := range doc.Fields {
		tmpMap[key] = field.Value
	}

	data, err := json.Marshal(tmpMap)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMarshalJSONFailed, err)
	}

	if err := json.Unmarshal(data, output); err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshalJSONFailed, err)
	}

	return nil
}
