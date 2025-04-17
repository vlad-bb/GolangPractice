package main

import (
	"GolangPractice/hw_06/document_store"
	"GolangPractice/hw_06/logger"
	"GolangPractice/hw_06/users"
	"fmt"
	"strconv"
)

const dumpFilePath = "store_dump.json.gz"

var log = logger.GetLogger()

func main() {
	store := document_store.NewStore()
	log.Info("New document_store created")

	userService := users.CreateService(*store, "id", "listUsers")

	names := []string{"Bob", "Mary", "John"}
	for i, name := range names {
		_, err := userService.CreateUser(strconv.Itoa(i), name)
		if err != nil {
			log.Error(err.Error())
		}
	}
	_, err := userService.CreateUser("5", "Vlad")

	user, err := userService.GetUser("1")
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug(fmt.Sprintf("User %s\n", user))

	userId := "1"
	err = userService.DeleteUser(userId)
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug(fmt.Sprintf("Deleted user %s\n", userId))

	listUsers, err := userService.ListUsers()
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug(fmt.Sprintf("List of users: %+v", listUsers))

	err = store.DumpToFile(dumpFilePath)
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Println("Dumped to file")
	storeFromDump, err := document_store.NewStoreFromFile(dumpFilePath)
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug(fmt.Sprintf("New document_store created from dump %v\n", storeFromDump))

	userServiceFromDump := users.CreateService(*storeFromDump, "id", "listUsers")
	_, err = userServiceFromDump.CreateUser("6", "Dmytro")
	listUsersFromDump, err := userServiceFromDump.ListUsers()
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug(fmt.Sprintf("List of users: %+v", listUsersFromDump))
}
