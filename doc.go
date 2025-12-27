package docstore

import (
	"math/rand"
	"time"
)

// DocId is a document identifier.
type DocId string

const EmptyDocId = DocId("")

// NewDocId creates a new document identifier.
func NewDocId(id string) DocId {
	return DocId(id)
}

// GenerateDocId creates a new random document identifier.
func GenerateDocId() DocId {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 12)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return DocId(b)
}

// String returns the string representation of the document identifier.
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
