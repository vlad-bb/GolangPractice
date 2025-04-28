package document_store

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
)

type Store struct {
	Storage map[string]*Collection `json:"storage"`
}

func NewStore() *Store {
	return &Store{
		Storage: make(map[string]*Collection),
	}

}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	_, ok := s.Storage[name]
	if ok {
		return false, nil
	}
	collection := NewCollection(*cfg)
	s.Storage[name] = collection
	return true, collection
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	collection, ok := s.Storage[name]
	return collection, ok
}

func (s *Store) DeleteCollection(name string) bool {
	_, ok := s.Storage[name]
	if ok {
		delete(s.Storage, name)

		return true
	}

	return false
}

func (s *Store) Dump() ([]byte, error) {
	// Методи повинен віддати dump нашого store, в який уключені дані про колекції та документ
	sBytes, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMarshalJSONFailed, err)
	}
	return sBytes, nil
}

func (s *Store) DumpToFile(filename string) error {
	jsonBytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrMarshalJSONFailed, err)
	}
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	if _, err := gzipWriter.Write(jsonBytes); err != nil {
		return fmt.Errorf("%w: %v", ErrCreateWriterFailed, err)
	}
	if err := gzipWriter.Close(); err != nil {
		return fmt.Errorf("%w: %v", ErrCloseFailed, err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("%w: %v", ErrWriteFailed, err)
	}
	return nil
}
func NewStoreFromDump(dump []byte) (*Store, error) {
	store := &Store{}
	err := json.Unmarshal(dump, store)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalJSONFailed, err)
	}
	return store, nil
}
func NewStoreFromFile(filename string) (*Store, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFileFailed, err)
	}

	reader, err := gzip.NewReader(bytes.NewReader(file))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateReaderFailed, err)
	}

	decompressedData := new(bytes.Buffer)
	_, err = decompressedData.ReadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadFromFailed, err)
	}

	if err := reader.Close(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCloseFailed, err)
	}

	store, err := NewStoreFromDump(decompressedData.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNewStoreFromDumpFailed, err)
	}
	return store, nil
}
