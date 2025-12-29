package docstore

import (
	"errors"
)

var (
	// ErrDocumentNotFound is returned when a document is not found.
	ErrDocumentNotFound = errors.New("document not found")

	// ErrEmptyDocumentId is returned when a valid document id is required.
	ErrEmptyDocumentId = errors.New("empty document id")

	// ErrDocumentTypeMismatch is returned when a document type does not match the expected type
	ErrDocumentTypeMismatch = errors.New("document type mismatch")
)
