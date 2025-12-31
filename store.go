package docstore

import (
	"iter"
	"sync"
)

var (
	store     map[DocId]Document
	storeLock sync.RWMutex
)

func init() {
	Clear()
}

// Clear removes all documents from the store.
func Clear() {
	storeLock.Lock()
	defer storeLock.Unlock()

	store = make(map[DocId]Document)
}

// Put adds a document to the store.
func Put(id DocId, doc Document) error {
	if id == EmptyDocId {
		return ErrEmptyDocumentId
	}

	storeLock.Lock()
	defer storeLock.Unlock()

	store[id] = doc
	return nil
}

// Get retrieves a document by Id
func Get(id DocId) (*Document, error) {
	if id == EmptyDocId {
		return nil, ErrEmptyDocumentId
	}

	storeLock.RLock()
	defer storeLock.RUnlock()

	doc, exists := store[id]
	if !exists {
		return nil, ErrDocumentNotFound
	}

	return &doc, nil
}

// GetAs retrieves a typed document by Id
func GetAs[T Document](id DocId) (*T, error) {
	doc, err := Get(id)
	if err != nil {
		return nil, err
	}

	typedDoc, ok := (*doc).(T)
	if !ok {
		return nil, ErrDocumentTypeMismatch
	}

	return &typedDoc, nil
}

// Delete removes a document
func Delete(id DocId) error {
	storeLock.Lock()
	defer storeLock.Unlock()

	// Check if document exists
	_, exists := store[id]
	if !exists {
		return ErrDocumentNotFound
	}

	// Remove document
	delete(store, id)

	return nil
}

// DeleteAllOf removes all documents of the given type
// returns the number of documents deleted
func DeleteAllOf[T Document]() int {
	storeLock.Lock()
	defer storeLock.Unlock()

	var marked []DocId

	for docId, doc := range store {
		if _, ok := doc.(T); ok {
			marked = append(marked, docId)
		}
	}

	count := 0

	for _, docId := range marked {
		delete(store, docId)
		count++
	}

	return count
}

// AllDocuments returns a iterator that goes over all stored documents
func AllDocuments() iter.Seq2[DocId, Document] {
	return func(yield func(DocId, Document) bool) {
		storeLock.RLock()
		defer storeLock.RUnlock()

		for docId, doc := range store {
			if !yield(docId, doc) {
				return
			}
		}
	}
}

// AllDocumentsOf returns a iterator that goes over all stored documents of a given type
func AllDocumentsOf[T Document]() iter.Seq2[DocId, T] {
	return func(yield func(DocId, T) bool) {
		storeLock.RLock()
		defer storeLock.RUnlock()

		for docId, doc := range store {
			typedDoc, ok := doc.(T)
			if ok && !yield(docId, typedDoc) {
				return
			}
		}
	}
}
