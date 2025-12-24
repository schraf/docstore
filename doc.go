package docstore

import (
	"math/rand"
	"time"
)

type DocId string

func NewDocId(id string) DocId {
	return DocId(id)
}

func GenerateDocId() DocId {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 12)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return DocId(b)
}

func (d DocId) String() string {
	return string(d)
}

// DocData used to restrict document data to concrete types
type DocData interface {
	any
}

// Document represents a stored document with metadata
type Document[T DocData] struct {
	Id        DocId     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Data      T         `json:"data"`
}
