package main

import (
	"GolangPractice/hw_08/lru"
	"fmt"
)

func main() {
	cache := lru.NewLruCache(2)
	fmt.Println("New LRU cache created")

	// Перевіряємо додавання і отримання значень
	cache.Put("lang", "Golang")
	value, ok := cache.Get("lang")
	if !ok {
		fmt.Printf("bad implementation '%s'", value)
	}
	fmt.Printf("Get value %s\n", value)

}
