package like_mongo

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrRecordParams   = errors.New("key and value cannot be empty")
)
