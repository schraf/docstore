package docstore

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
)

func hashSorter[T DocData](left *Document[T], right *Document[T]) bool {
	return left.Id.String() < right.Id.String()
}

// Hash returns a hash of the documents in the store.
func (s *Store[T]) Hash() (int32, error) {
	var documents []*Document[T]

	for _, doc := range s.documents {
		documents = append(documents, &doc)
	}

	documents = sortDocuments(documents, hashSorter)

	// Create a new FNV-1a hash
	h := fnv.New32a()

	// Define a struct for hashing that excludes timestamps
	type hashableDoc struct {
		Id   DocId `json:"id"`
		Data T     `json:"data"`
	}

	// Marshal and write each document to the hash
	for _, doc := range documents {
		b, err := json.Marshal(hashableDoc{Id: doc.Id, Data: doc.Data})
		if err != nil {
			return 0, fmt.Errorf("failed to marshal document %s for hashing: %w", doc.Id, err)
		}
		if _, err := h.Write(b); err != nil {
			return 0, fmt.Errorf("failed to write document %s to hash: %w", doc.Id, err)
		}
	}

	return int32(h.Sum32()), nil
}
