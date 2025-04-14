package users

import (
	"GolangPractice/hw_05/document_store"
	"fmt"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	coll document_store.Collection
}

func CreateService(collection document_store.Collection) *Service {
	return &Service{coll: collection}
}

func (s *Service) CreateUser(payload map[string]interface{}) (*User, error) {
	doc, err := document_store.MarshalDocument(payload)
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = document_store.UnmarshalDocument(doc, user)
	if err != nil {
		return nil, err
	}
	err = s.coll.Put(*doc)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) ListUsers() ([]User, error) {
	docList := s.coll.List()
	users := make([]User, len(docList))
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
