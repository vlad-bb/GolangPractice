package document_store

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
)

type Store struct {
	storage map[string]*Collection
}

func NewStore() *Store {
	return &Store{
		storage: make(map[string]*Collection),
	}

}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	_, ok := s.storage[name]
	if ok {
		return false, nil
	}
	collection := NewCollection(*cfg)
	s.storage[name] = collection
	return true, collection
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	collection, ok := s.storage[name]
	return collection, ok
}

func (s *Store) DeleteCollection(name string) bool {
	_, ok := s.storage[name]
	if ok {
		delete(s.storage, name)

		return true
	}

	return false
}

func (s *Store) Dump() ([]byte, error) {
	// Методи повинен віддати дамп нашого стору в який включені дані про колекції та документ
	sBytes, err := json.Marshal(&struct {
		Storage map[string]*Collection `json:"storage"`
	}{
		Storage: s.storage,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMarshalJSONFailed, err)
	}
	return sBytes, nil
}

func (s *Store) DumpToFile(filename string) error {
	jsonBytes, err := json.MarshalIndent(&struct {
		Storage map[string]*Collection `json:"storage"`
	}{
		Storage: s.storage,
	}, "", "    ")
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
	aux := &struct {
		Storage map[string]*Collection `json:"storage"`
	}{}
	err := json.Unmarshal(dump, aux)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalJSONFailed, err)
	}
	store.storage = aux.Storage
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
