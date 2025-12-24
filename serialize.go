package docstore

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"
)

const (
	snapshotMagic = "DSS1"
)

type snapshotInfo struct {
	Magic     string    `json:"magic"`
	Timestamp time.Time `json:"timestamp"`
	DocType   string    `json:"doc_type"`
	Hash      int32     `json:"hash"`
}

func (s *Store[T]) WriteTo(w io.Writer) (int64, error) {
	hash, err := s.Hash()
	if err != nil {
		return 0, fmt.Errorf("unable to hash documents for snapshot: %w", err)
	}

	typeName := reflect.TypeFor[T]().String()

	info := snapshotInfo{
		Magic:     snapshotMagic,
		Timestamp: time.Now().UTC(),
		DocType:   typeName,
		Hash:      hash,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(info); err != nil {
		return 0, fmt.Errorf("failed to encode and write snapshot info: %w", err)
	}

	// Write out each document
	for _, doc := range s.documents {
		if err := encoder.Encode(doc); err != nil {
			return 0, fmt.Errorf("failed to encode document %s: %w", doc.Id, err)
		}
	}

	return 0, nil
}

func (s *Store[T]) ReadFrom(r io.Reader) (int64, error) {
	typeName := reflect.TypeFor[T]().String()
	decoder := json.NewDecoder(r)

	// Read out the snapshot info
	info := snapshotInfo{}
	if err := decoder.Decode(&info); err != nil {
		return 0, fmt.Errorf("failed to decode snapshot info: %w", err)
	}

	if info.Magic != snapshotMagic {
		return 0, ErrInvalidSnapshotMagic
	}

	if info.DocType != typeName {
		return 0, ErrMismatchedDocType
	}

	// Read out each document
	for decoder.More() {
		var doc Document[T]

		if err := decoder.Decode(&doc); err != nil {
			return 0, fmt.Errorf("failed to decode document: %w", err)
		}

		s.documents[doc.Id] = doc
	}

	return 0, nil
}
