package docstore

import "math/rand"

const EmptyDocId = DocId("")

// DocId is a document identifier.
type DocId string

// NewDocId creates a new document identifier.
func NewDocId(id string) DocId {
	return DocId(id)
}

// String returns the string representation of the document identifier.
func (d DocId) String() string {
	return string(d)
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

type Document any
