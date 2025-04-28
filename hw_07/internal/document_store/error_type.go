package document_store

import "errors"

var (
	ErrCollectionAlreadyExists  = errors.New("collection already exists")
	ErrUnsupportedDocumentField = errors.New("unsupported document field")
	ErrDocumentNil              = errors.New("document is nil")
	ErrUnmarshalJSONFailed      = errors.New("json unmarshal to struct failed")
	ErrMarshalJSONFailed        = errors.New("json marshal failed")
	ErrUnsupportedFieldType     = errors.New("unsupported field type")
	ErrPrimaryKeyNotFound       = errors.New("primary key field not found")
	ErrPrimaryKeyWrongType      = errors.New("primary key field is not of type string")
	ErrPrimaryKeyInvalidValue   = errors.New("primary key value is invalid")
	ErrWriteFailed              = errors.New("write failed")
	ErrCloseFailed              = errors.New("close failed")
	ErrCreateWriterFailed       = errors.New("create writer failed")
	ErrCreateReaderFailed       = errors.New("failed to create gzip reader")
	ErrReadFileFailed           = errors.New("failed to read file")
	ErrReadFromFailed           = errors.New("failed to read from file")
	ErrNewStoreFromDumpFailed   = errors.New("failed to create new store from dump")
)
