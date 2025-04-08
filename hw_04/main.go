package main

import (
	"GolangPractice/hw_04/documentstore"
	"fmt"
)

func main() {
	store := documentstore.NewStore() // create store
	fmt.Printf("New store created %v\n", store)

	cfg := &documentstore.CollectionConfig{
		PrimaryKey: "pk",
	}

	// Test Collection CRUD
	ok, usersCollection := store.CreateCollection("users", cfg) // success case
	if ok {
		fmt.Printf("Created collection %v\n", usersCollection)
	} else {
		fmt.Println("Collection already exists")
	}

	ok, usersCollection2 := store.CreateCollection("users", cfg) // failed case
	if ok {
		fmt.Printf("Created collection %v\n", usersCollection2)
	} else {
		fmt.Println("Collection already exists")
	}

	usersCollection, ok = store.GetCollection("users") // success case
	if ok {
		fmt.Printf("Get collection %v\n", usersCollection)
	} else {
		fmt.Println("Collection doesn't exists")
	}

	usersCollection2, ok = store.GetCollection("users2") // failed case
	if ok {
		fmt.Printf("Get collection %v\n", usersCollection2)
	} else {
		fmt.Println("Collection doesn't exists")
	}

	ok = store.DeleteCollection("users") // success case
	if ok {
		fmt.Println("Collection deleted successfully")
	} else {
		fmt.Println("Collection doesn't exists")
	}

	ok = store.DeleteCollection("users2") // success case
	if ok {
		fmt.Println("Collection deleted successfully")
	} else {
		fmt.Println("Collection doesn't exists")
	}

	// Test Document CRUD
	ok, usersCollection = store.CreateCollection("users", cfg) // success case
	if ok {
		fmt.Printf("Created collection %v\n", usersCollection)
	} else {
		fmt.Println("Collection already exists")
	}

	user1 := documentstore.Document{Fields: map[string]documentstore.DocumentField{
		"pk": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "1",
		},
		"user_name": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "John Doe",
		},
		"age": {
			Type:  documentstore.DocumentFieldTypeNumber,
			Value: 42,
		},
		"is_active": {
			Type:  documentstore.DocumentFieldTypeBool,
			Value: false,
		},
	}}

	user2 := documentstore.Document{Fields: map[string]documentstore.DocumentField{
		"pk": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "2",
		},
		"user_name": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "Bob Black",
		},
		"age": {
			Type:  documentstore.DocumentFieldTypeNumber,
			Value: 18,
		},
		"is_active": {
			Type:  documentstore.DocumentFieldTypeBool,
			Value: true,
		},
	}}

	user3 := documentstore.Document{Fields: map[string]documentstore.DocumentField{
		"user_name": {
			Type:  documentstore.DocumentFieldTypeString,
			Value: "Alan White",
		},
		"age": {
			Type:  documentstore.DocumentFieldTypeNumber,
			Value: 20,
		},
		"is_active": {
			Type:  documentstore.DocumentFieldTypeBool,
			Value: true,
		},
	}}

	usersCollection.Put(user1) // success case
	usersCollection.Put(user2) // success case
	usersCollection.Put(user3) // fail test

	allUsers := usersCollection.List()
	fmt.Println(allUsers) // only 2 elements

	user, ok := usersCollection.Get("1") // success case
	if ok {
		fmt.Printf("User found %v\n", user)
	} else {
		fmt.Println("User not found")
	}

	user, ok = usersCollection.Get("3") // failed case
	if ok {
		fmt.Printf("User found %v\n", user)
	} else {
		fmt.Println("User not found")
	}

	if usersCollection.Delete("1") { // success case
		fmt.Println("User deleted")
	} else {
		fmt.Println("User not found")
	}

	if usersCollection.Delete("3") { // failed case
		fmt.Println("User deleted")
	} else {
		fmt.Println("User not found")
	}
}
