package main

import (
	"GolangPractice/hw_11/internal/document_store"
	"GolangPractice/hw_11/internal/llogger"
	"GolangPractice/hw_11/internal/users"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"strconv"
	"sync"
)

var logger = llogger.SetupLogger()

var store = document_store.NewStore()
var userService = users.CreateService(store, "id", "listUsers")
var iterationSteps = 10

func main() {
	wg := sync.WaitGroup{}
	gofakeit.Seed(0)
	for i := 0; i < iterationSteps; i++ {
		wg.Add(1)
		name := gofakeit.Name()
		id := strconv.Itoa(i)
		go userCrud(id, name, &wg)
	}
	wg.Wait()

}

func userCrud(id string, name string, wg *sync.WaitGroup) {
	logger.Debug(fmt.Sprintf("Start CRUD for user: %s %s", id, name))
	defer wg.Done()
	_, err := userService.CreateUser(id, name)
	if err != nil {
		logger.Error(err.Error())
	}
	user, err := userService.GetUser(id)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Debug(fmt.Sprintf("User %s\n", user))

	err = userService.DeleteUser(id)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Debug(fmt.Sprintf("Deleted user %s\n", id))

	listUsers, err := userService.ListUsers()
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Debug(fmt.Sprintf("List of users: %+v", listUsers))
}
