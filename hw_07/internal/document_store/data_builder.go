package document_store

type DocumentBuilder struct {
	fields map[string]DocumentField
}

func NewDocumentBuilder() *DocumentBuilder {
	return &DocumentBuilder{
		fields: make(map[string]DocumentField),
	}
}

func (b *DocumentBuilder) WithField(name string, fieldType DocumentFieldType, value interface{}) *DocumentBuilder {
	b.fields[name] = DocumentField{
		Type:  fieldType,
		Value: value,
	}
	return b
}

func (b *DocumentBuilder) Build() Document {
	return Document{Fields: b.fields}
}
