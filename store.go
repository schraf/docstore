package docstore

import (
	"time"
)

// Store is a simple document store.
type Store[T DocData] struct {
	documents map[DocId]Document[T]
}

// NewStore creates a new store.
func NewStore[T DocData]() *Store[T] {
	return &Store[T]{
		documents: make(map[DocId]Document[T]),
	}
}

// Clear all documents
func (s *Store[T]) Clear() error {
	s.documents = make(map[DocId]Document[T])
	return nil
}

// Put adds or updates a document
func (s *Store[T]) Put(doc Document[T]) error {
	// Check if document already exists
	if oldDoc, exists := s.documents[doc.Id]; exists {
		// Preserve creation time, update modification time
		doc.CreatedAt = oldDoc.CreatedAt
		doc.UpdatedAt = time.Now()

		// Update document
		s.documents[doc.Id] = doc
		return nil
	}

	// Set timestamps
	now := time.Now()
	doc.CreatedAt = now
	doc.UpdatedAt = now

	// Store document
	s.documents[doc.Id] = doc

	return nil
}

// Get retrieves a document by Id
func (s *Store[T]) Get(id DocId) (*Document[T], error) {
	doc, exists := s.documents[id]
	if !exists {
		return nil, ErrDocumentNotFound
	}

	return &doc, nil
}

// Delete removes a document
func (s *Store[T]) Delete(id DocId) error {
	// Check if document exists
	_, exists := s.documents[id]
	if !exists {
		return ErrDocumentNotFound
	}

	// Remove document
	delete(s.documents, id)

	return nil
}
