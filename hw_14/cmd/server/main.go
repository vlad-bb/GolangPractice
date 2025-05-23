package main

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	opts := options.Client()

	opts.ApplyURI(DbUri)

	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("error connecting to MongoDB: %v", err))
	}

	handler := NewUserHandler(c)

	http.HandleFunc("POST /put_document", handler.handlePut)
	http.HandleFunc("POST /get_document", handler.handleGet)
	http.HandleFunc("POST /delete_document", handler.handleDelete)
	http.HandleFunc("POST /list_documents", handler.handleList)
	http.HandleFunc("POST /create_collection", handler.handlePutCollection)
	http.HandleFunc("POST /list_collections", handler.handleListCollection)
	http.HandleFunc("POST /delete_collection", handler.handleDeleteCollection)
	http.HandleFunc("POST /create_index", handler.handlePutIndex)
	http.HandleFunc("POST /delete_index", handler.handleDeleteIndex)

	logger.Debug("server listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(fmt.Errorf("server listening failed: %v", err))
	}

}
