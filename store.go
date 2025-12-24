package docstore

import (
	"time"
)

type Store[T DocData] struct {
	documents map[DocId]Document[T]
}

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

// Select finds documents based on criteria
func (s *Store[T]) Select(query Query[T]) (*QueryResult[T], error) {
	// Collect all documents for this schema
	var candidates []*Document[T]
	for _, doc := range s.documents {
		if s.filter(&doc, query.Filters) {
			candidates = append(candidates, &doc)
		}
	}

	if query.SortBy != nil {
		candidates = sortDocuments(candidates, *query.SortBy)
	}

	// Calculate pagination
	total := len(candidates)

	// Apply limit
	if query.Limit > 0 && len(candidates) > query.Limit {
		candidates = candidates[:query.Limit]
	}

	return &QueryResult[T]{
		Documents: candidates,
		Total:     total,
	}, nil
}

// Helper methods

func (s *Store[T]) filter(doc *Document[T], filters []QueryFilter[T]) bool {
	if filters == nil {
		return true
	}

	for _, filter := range filters {
		if !filter(doc) {
			return false
		}
	}

	return true
}
