package like_mongo

import (
	"GolangPractice/hw_12/internal/document_store"
	"GolangPractice/hw_12/internal/llogger"
	"fmt"
)

var logger = llogger.SetupLogger()

type Record struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Service struct {
	coll *document_store.Collection
}

func CreateService(store *document_store.Store, pk string, collectionName string) Service {
	cfg := &document_store.CollectionConfig{
		PrimaryKey: pk,
	}
	ok, collection := store.CreateCollection(collectionName, cfg)
	if !ok {
		logger.Info(fmt.Sprintf("Msg: %v", document_store.ErrCollectionAlreadyExists))
		collection, _ = store.GetCollection(collectionName)
	}
	logger.Info("Service created")
	return Service{coll: collection}
}

func (s *Service) CreateRecord(key string, value string) (*Record, error) {
	if key == "" || value == "" {
		return nil, ErrRecordParams
	}
	doc := &document_store.Document{
		Fields: make(map[string]document_store.DocumentField),
	}
	doc.Fields["key"] = document_store.DocumentField{
		Type:  document_store.DocumentFieldTypeString,
		Value: key,
	}
	doc.Fields["value"] = document_store.DocumentField{
		Type:  document_store.DocumentFieldTypeString,
		Value: value,
	}
	err := s.coll.Put(*doc)
	if err != nil {
		return nil, err
	}
	record := &Record{
		Key:   key,
		Value: value,
	}
	logger.Info("Record created")
	return record, nil
}

func (s *Service) ListRecords() ([]Record, error) {
	docList := s.coll.List()
	records := make([]Record, 0, len(docList))
	for _, doc := range docList {
		record := &Record{}
		err := document_store.UnmarshalDocument(&doc, record)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}
	logger.Info("Records listed")
	return records, nil
}

func (s *Service) GetRecord(key string) (*Record, error) {
	doc, ok := s.coll.Get(key)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrRecordNotFound, key)
	}
	record := &Record{}
	err := document_store.UnmarshalDocument(doc, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *Service) DeleteRecord(key string) error {
	if ok := s.coll.Delete(key); !ok {
		return fmt.Errorf("%w: %s", ErrRecordNotFound, key)
	}
	logger.Info("Record deleted")
	return nil
}
