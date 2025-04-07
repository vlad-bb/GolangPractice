package documentstore

import "fmt"

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

var documents = map[string]Document{}

func Put(doc Document) {
	// 1. Перевірити що документ містить в мапі поле `key` типу `string`
	// 2. Додати Document до локальної мапи з документами
	field, err := doc.Fields["key"]
	if err != true {
		fmt.Printf("Not any keys")
		return
	}

	if field.Type != DocumentFieldTypeString {
		fmt.Printf("Invalid  key %v, should be 'string'\n", field.Type)
		return
	}
	key, ok := field.Value.(string)
	if !ok {
		fmt.Println("field type is not a string")
		return
	}
	if key == "" {
		fmt.Println("that key shouldn't be empty")
		return
	}

	documents[key] = doc
}

func Get(key string) (*Document, bool) {
	// Потрібно повернути документ по ключу
	// Якщо документ знайдено, повертаємо `true` та поінтер на документ
	// Інакше повертаємо `false` та `nil`
	doc, ok := documents[key]
	if !ok {
		return nil, false
	}
	return &doc, true
}

func Delete(key string) bool {
	// Видаляємо документа по ключу.
	// Повертаємо `true` якщо ми знайшли і видалили документові
	// Повертаємо `false` якщо документ не знайдено
	_, ok := Get(key)
	if ok {
		delete(documents, key)
		return true
	} else {
		return false
	}
}

func List() []Document {
	// Повертаємо список усіх документів
	docs := make([]Document, 0, len(documents))
	for _, doc := range documents {
		docs = append(docs, doc)
	}
	return docs
}
