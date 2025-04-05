package main

import (
	"fmt"
	"hw_03/documentstore"
)

var firstDoc = documentstore.Document{
	Fields: map[string]documentstore.DocumentField{
		"key":    {Type: documentstore.DocumentFieldTypeString, Value: "first"},
		"amount": {Type: documentstore.DocumentFieldTypeNumber, Value: 100},
	},
}

var thirdDoc = documentstore.Document{
	Fields: map[string]documentstore.DocumentField{
		"key":    {Type: documentstore.DocumentFieldTypeString, Value: "third"},
		"amount": {Type: documentstore.DocumentFieldTypeNumber, Value: 100},
	},
}

var fourthDoc = documentstore.Document{
	Fields: map[string]documentstore.DocumentField{
		"key":    {Type: documentstore.DocumentFieldTypeString, Value: "fourth"},
		"amount": {Type: documentstore.DocumentFieldTypeNumber, Value: 100},
	},
}

func main() {
	documentstore.Put(firstDoc)  // add document
	documentstore.Put(thirdDoc)  // add document
	documentstore.Put(fourthDoc) // add document

	doc, ok := documentstore.Get("first") // success case
	if ok {
		fmt.Printf("Document: %v\n", doc)
	} else {
		fmt.Println("Document not found")
	}
	doc, ok = documentstore.Get("second") // fail case
	if ok {
		fmt.Printf("Document: %v\n", doc)
	} else {
		fmt.Println("Document not found")
	}

	status := documentstore.Delete("first") // success case
	if status {
		fmt.Println("Document deleted")
	} else {
		fmt.Println("Document not found")
	}
	status = documentstore.Delete("second") // fail case
	if status {
		fmt.Println("Document deleted")
	} else {
		fmt.Println("Document not found")
	}

	docs := documentstore.List()
	fmt.Println(docs)
}
