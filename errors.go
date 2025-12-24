package docstore

import (
	"errors"
)

var (
	ErrDocumentNotFound     = errors.New("document not found")
	ErrInvalidSnapshotMagic = errors.New("invalid snapshot magic identifier")
	ErrMismatchedDocType    = errors.New("snapshot doc type does not match store doc type")
)
