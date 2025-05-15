package commands

import "GolangPractice/hw_13/internal/like_mongo"

type PutCommandRequestPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PutCommandResponsePayload struct{}

type GetCommandRequestPayload struct {
	Key string `json:"key"`
}

type GetCommandResponsePayload struct {
	Value string `json:"value"`
	Ok    bool   `json:"ok"`
}

type DeleteCommandRequestPayload struct {
	Key string `json:"key"`
}

type DeleteCommandResponsePayload struct {
	Ok bool `json:"ok"`
}

type ExitCommandResponsePayload struct {
	Msg string `json:"msg"`
}

type ListCommandResponsePayload struct {
	Items []like_mongo.Record `json:"items"`
}

const (
	PutCommandName    string = "put"
	GetCommandName    string = "get"
	DeleteCommandName string = "delete"
	ListCommandName   string = "list"
)
