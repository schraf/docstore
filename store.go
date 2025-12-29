package docstore

import "sync"

var (
	store     map[DocId]Document
	storeLock sync.RWMutex
)

func init() {
	Clear()
}

func Clear() {
	storeLock.Lock()
	defer storeLock.Unlock()

	store = make(map[DocId]Document)
}

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
