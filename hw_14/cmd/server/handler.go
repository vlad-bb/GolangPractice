package main

import (
	"encoding/json"
	"fmt"
	"hw_14/internal/llogger"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logger = llogger.SetupLogger()

type Handler struct {
	coll *mongo.Collection
	db   *mongo.Database
}

func NewUserHandler(c *mongo.Client) *Handler {
	logger.Debug("Create new UserHandler")
	db := c.Database(dbName)
	coll := db.Collection(collectionName)
	return &Handler{coll: coll, db: db}
}

type DocumentBody struct {
	UserId   string `json:"user_id" bson:"user_id"`
	Name     string `json:"name" bson:"name"`
	Age      int    `json:"age" bson:"age"`
	IsActive bool   `json:"is_active" bson:"is_active"`
}

type PutReqBody struct {
	CollectionName string       `json:"collection_name"`
	Document       DocumentBody `json:"document"`
}

type PutRespBody struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handlePut(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handle put request started")
	reqBody := PutReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	currentColl := h.coll
	if collectionName != reqBody.CollectionName {
		currentColl = h.db.Collection(reqBody.CollectionName)
	}
	doc := &DocumentBody{
		UserId:   reqBody.Document.UserId,
		Name:     reqBody.Document.Name,
		Age:      reqBody.Document.Age,
		IsActive: reqBody.Document.IsActive,
	}
	filter := bson.M{"user_id": doc.UserId}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err = currentColl.UpdateOne(r.Context(), filter, update, opts)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to update document: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := PutRespBody{Ok: true}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
	logger.Debug("handle put request finished")
}

type GetReqBody struct {
	UserId         string `json:"user_id"`
	CollectionName string `json:"collection_name"`
}

type GetRespBody struct {
	Document DocumentBody `json:"document"`
	Ok       bool         `json:"ok"`
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handle get request started")
	reqBody := GetReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	logger.Debug(reqBody.UserId)
	filter := bson.M{"user_id": reqBody.UserId}
	doc := &DocumentBody{}
	currentColl := h.coll
	if collectionName != reqBody.CollectionName {
		currentColl = h.db.Collection(reqBody.CollectionName)
	}
	err = currentColl.FindOne(r.Context(), filter).Decode(doc)
	if err != nil {
		http.Error(w, fmt.Errorf("not found document: %w", err).Error(), http.StatusNotFound)
		return
	}

	respBody := GetRespBody{Document: *doc, Ok: true}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type DeleteReqBody struct {
	UserId         string `json:"user_id"`
	CollectionName string `json:"collection_name"`
}

type DeleteRespBody struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	reqBody := DeleteReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	var res map[string]interface{}
	filter := bson.M{"user_id": reqBody.UserId}
	currentColl := h.coll
	if collectionName != reqBody.CollectionName {
		currentColl = h.db.Collection(reqBody.CollectionName)
	}
	err = currentColl.FindOneAndDelete(r.Context(), filter).Decode(&res)
	if err != nil {
		http.Error(w, fmt.Errorf("not found document for delete: %w", err).Error(), http.StatusNotFound)
		return
	}

	respBody := DeleteRespBody{Ok: true}
	err = json.NewEncoder(w).Encode(respBody)

	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type ListReqBody struct {
	CollectionName string `json:"collection_name"`
}

type ListRespBody struct {
	Docs []DocumentBody `json:"documents"`
}

func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	reqBody := ListReqBody{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	currentColl := h.coll
	if collectionName != reqBody.CollectionName {
		currentColl = h.db.Collection(reqBody.CollectionName)
	}
	cur, err := currentColl.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, fmt.Errorf("failed to find documents: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	var docs []DocumentBody
	err = cur.All(r.Context(), &docs)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode documents: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := ListRespBody{Docs: docs}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type GetReqBodyColl struct {
	CollectionName string `json:"collection_name"`
}

type GetRespBodyColl struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handlePutCollection(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handle get collection request started")
	reqBody := GetReqBodyColl{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode req body: %w", err).Error(), http.StatusBadRequest)
		return
	}
	logger.Debug(reqBody.CollectionName)
	h.db.Collection(reqBody.CollectionName)

	respBody := GetRespBodyColl{Ok: true}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode resp body: %w", err).Error(), http.StatusInternalServerError)
	}
}

type CollectionBody struct {
	CollectionName string `json:"collection_name"`
}

type ListRespBodyColl struct {
	Collection []CollectionBody `json:"collections"`
	Ok         bool             `json:"ok"`
}

func (h *Handler) handleListCollection(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handle list collection request started")

	collections, err := h.db.ListCollectionNames(r.Context(), bson.D{})
	if err != nil {
		http.Error(w, fmt.Errorf("failed to list collections: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	var collectionList []CollectionBody
	for _, name := range collections {
		collectionList = append(collectionList, CollectionBody{CollectionName: name})
	}

	respBody := ListRespBodyColl{
		Collection: collectionList,
		Ok:         true,
	}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode response body: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug("handle list collection request finished")
}

type DeleteReqBodyColl struct {
	CollectionName string `json:"collection_name"`
}

type DeleteRespBodyColl struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handleDeleteCollection(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handle delete collection request started")

	reqBody := DeleteReqBodyColl{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to decode request body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if reqBody.CollectionName == "" {
		http.Error(w, "collection_name is required", http.StatusBadRequest)
		return
	}

	err = h.db.Collection(reqBody.CollectionName).Drop(r.Context())
	if err != nil {
		http.Error(w, fmt.Errorf("failed to delete collection: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := DeleteRespBodyColl{Ok: true}
	err = json.NewEncoder(w).Encode(respBody)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to encode response body: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug("handle delete collection request finished")
}

type CreateIndexReqBody struct {
	CollectionName string `json:"collection_name"`
	IndexField     string `json:"index_field"`
}

type CreateIndexRespBody struct {
	Ok        bool   `json:"ok"`
	IndexName string `json:"index_name"`
}

func (h *Handler) handlePutIndex(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handle create index request started")

	var reqBody CreateIndexReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, fmt.Errorf("failed to decode request body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if reqBody.CollectionName == "" || reqBody.IndexField == "" {
		http.Error(w, "collection_name and index_field are required", http.StatusBadRequest)
		return
	}

	collection := h.db.Collection(reqBody.CollectionName)

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: reqBody.IndexField, Value: 1}},
	}

	indexName, err := collection.Indexes().CreateOne(r.Context(), indexModel)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to create index: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	respBody := CreateIndexRespBody{
		Ok:        true,
		IndexName: indexName,
	}

	if err := json.NewEncoder(w).Encode(respBody); err != nil {
		http.Error(w, fmt.Errorf("failed to encode response: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug("handle create index request finished")
}

type DeleteIndexReqBody struct {
	CollectionName string `json:"collection_name"`
	IndexField     string `json:"index_field"`
}

type DeleteIndexRespBody struct {
	Ok bool `json:"ok"`
}

func (h *Handler) handleDeleteIndex(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handle delete index request started")

	var reqBody DeleteIndexReqBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, fmt.Errorf("failed to decode request body: %w", err).Error(), http.StatusBadRequest)
		return
	}

	if reqBody.CollectionName == "" || reqBody.IndexField == "" {
		http.Error(w, "collection_name and index_field are required", http.StatusBadRequest)
		return
	}

	collection := h.db.Collection(reqBody.CollectionName)

	indexName := fmt.Sprintf("%s_1", reqBody.IndexField)

	_, err := collection.Indexes().DropOne(r.Context(), indexName)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to delete index: %w", err).Error(), http.StatusNotFound)
		return
	}

	respBody := DeleteIndexRespBody{Ok: true}
	if err := json.NewEncoder(w).Encode(respBody); err != nil {
		http.Error(w, fmt.Errorf("failed to encode response: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug("handle delete index request finished")
}
