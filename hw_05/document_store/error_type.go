package document_store

import "errors"

var (
	ErrDocumentNotFound         = errors.New("document not found")
	ErrCollectionAlreadyExists  = errors.New("collection already exists")
	ErrCollectionNotFound       = errors.New("collection not found")
	ErrUnsupportedDocumentField = errors.New("unsupported document field")
	ErrDocumentNil              = errors.New("document is nil")
	ErrUnmarshalJSONFailed      = errors.New("json unmarshal to struct failed")
	ErrMarshalJSONFailed        = errors.New("json marshal from map failed")
	ErrUnsupportedFieldType     = errors.New("unsupported field type")
	ErrPrimaryKeyNotFound       = errors.New("primary key field not found")
	ErrPrimaryKeyWrongType      = errors.New("primary key field is not of type string")
	ErrPrimaryKeyInvalidValue   = errors.New("primary key value is invalid")
)
