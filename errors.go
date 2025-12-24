package docstore

import (
	"errors"
)

var (
	// ErrDocumentNotFound is returned when a document is not found.
	ErrDocumentNotFound = errors.New("document not found")
)
