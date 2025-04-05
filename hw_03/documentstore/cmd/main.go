package main

import (
	"hw_03/documentstore"
)

var firstDoc = documentstore.Document{
	Fields: map[string]documentstore.DocumentField{
		"key":    {Type: documentstore.DocumentFieldTypeString, Value: "first"},
		"amount": {Type: documentstore.DocumentFieldTypeNumber, Value: 100},
	},
}

func main() {
	documentstore.Put(firstDoc) // Додаємо документ

	//fmt.Println(documentstore.documents["first"]) // так не можно добратися до змінної бо вона не доступна бо з малої літери
}
