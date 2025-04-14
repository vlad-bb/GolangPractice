package main

import (
	"GolangPractice/hw_05/document_store"
	"GolangPractice/hw_05/users"
	"fmt"
	"strconv"
)

func main() {
	store := document_store.NewStore()
	fmt.Printf("New document_store created %v\n", store)

	cfg := &document_store.CollectionConfig{
		PrimaryKey: "id",
	}

	ok, usersCollection := store.CreateCollection("listUsers", cfg)
	if !ok {
		fmt.Printf("failed to create collection: %v", document_store.ErrCollectionAlreadyExists)
	} else {
		usersCollection, _ = store.GetCollection("listUsers")
	}

	userService := users.CreateService(*usersCollection)

	names := []string{"Bob", "Mary", "John"}
	for i, name := range names {
		payload := map[string]interface{}{"id": strconv.Itoa(i), "name": name}
		_, err := userService.CreateUser(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
	payload := map[string]interface{}{"id": "5", "name": "Vlad"}
	_, err := userService.CreateUser(payload)

	user, err := userService.GetUser("1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)

	userId := "1"
	err = userService.DeleteUser(userId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Deleted user %s\n", userId)

	listUsers, err := userService.ListUsers()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(listUsers)
}
