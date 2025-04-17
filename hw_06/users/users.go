package users

import (
	"GolangPractice/hw_06/document_store"
	"GolangPractice/hw_06/logger"
	"fmt"
)

var log = logger.GetLogger()

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	coll *document_store.Collection
}

func CreateService(store document_store.Store, pk string, collectionName string) Service {
	cfg := &document_store.CollectionConfig{
		PrimaryKey: pk,
	}
	ok, collection := store.CreateCollection(collectionName, cfg)
	if !ok {
		log.Info(fmt.Sprintf("Msg: %v", document_store.ErrCollectionAlreadyExists))
		collection, _ = store.GetCollection(collectionName)
	}
	return Service{coll: collection}
}

func (s *Service) CreateUser(newGuid string, name string) (*User, error) {
	if newGuid == "" || name == "" {
		return nil, ErrUserParams
	}
	doc := &document_store.Document{
		Fields: make(map[string]document_store.DocumentField),
	}
	doc.Fields["id"] = document_store.DocumentField{
		Type:  document_store.DocumentFieldTypeString,
		Value: newGuid,
	}
	doc.Fields["name"] = document_store.DocumentField{
		Type:  document_store.DocumentFieldTypeString,
		Value: name,
	}
	err := s.coll.Put(*doc)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:   newGuid,
		Name: name,
	}
	return user, nil
}

func (s *Service) ListUsers() ([]User, error) {
	docList := s.coll.List()
	users := make([]User, 0, len(docList))
	for _, doc := range docList {
		user := &User{}
		err := document_store.UnmarshalDocument(&doc, user)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

func (s *Service) GetUser(userID string) (*User, error) {
	doc, ok := s.coll.Get(userID)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUserNotFound, userID)
	}
	user := &User{}
	err := document_store.UnmarshalDocument(doc, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) DeleteUser(userID string) error {
	if ok := s.coll.Delete(userID); !ok {
		return fmt.Errorf("%w: %s", ErrUserNotFound, userID)
	}
	return nil
}
