package docstore

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SerializeTestDoc struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSerialize_Success(t *testing.T) {
	// Setup original store
	store := NewStore[SerializeTestDoc]()
	docs := []Document[SerializeTestDoc]{
		{Id: GenerateDocId(), Data: SerializeTestDoc{Name: "John Doe", Age: 30}},
		{Id: GenerateDocId(), Data: SerializeTestDoc{Name: "Jane Smith", Age: 25}},
	}
	for _, doc := range docs {
		err := store.Put(doc)
		require.NoError(t, err)
	}

	originalHash, err := store.Hash()
	require.NoError(t, err)

	// Write snapshot to buffer
	var buf bytes.Buffer
	_, err = store.WriteTo(&buf)
	require.NoError(t, err)

	// Read from snapshot into a new store
	newStore := NewStore[SerializeTestDoc]()
	_, err = newStore.ReadFrom(&buf)
	require.NoError(t, err)

	// Verify new store
	newHash, err := newStore.Hash()
	require.NoError(t, err)
	assert.Equal(t, originalHash, newHash)

	retrievedDoc, err := newStore.Get(docs[0].Id)
	require.NoError(t, err)
	assert.Equal(t, docs[0].Data.Name, retrievedDoc.Data.Name)
}

func TestSerialize_InvalidMagic(t *testing.T) {
	badSnapshot := strings.NewReader(`{
		"magic": "BADMAGIC",
		"timestamp": "2025-01-01T00:00:00Z",
		"doc_type": "docstore.SerializeTestDoc",
		"hash": 12345
	}`)

	store := NewStore[SerializeTestDoc]()
	_, err := store.ReadFrom(badSnapshot)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidSnapshotMagic)
}

func TestSerialize_MismatchedDocType(t *testing.T) {
	type DifferentDoc struct {
		Value float64 `json:"value"`
	}

	snapshot := strings.NewReader(`{
		"magic": "DSS1",
		"timestamp": "2025-01-01T00:00:00Z",
		"doc_type": "docstore.SerializeTestDoc",
		"hash": 12345
	}`)

	store := NewStore[DifferentDoc]()
	_, err := store.ReadFrom(snapshot)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMismatchedDocType)
}
