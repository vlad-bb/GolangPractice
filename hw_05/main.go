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

	userService := users.CreateService(*store, "id", "listUsers")

	names := []string{"Bob", "Mary", "John"}
	for i, name := range names {
		_, err := userService.CreateUser(strconv.Itoa(i), name)
		if err != nil {
			fmt.Println(err)
		}
	}
	_, err := userService.CreateUser("5", "Vlad")

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
