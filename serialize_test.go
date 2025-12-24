package docstore

import (
	"bytes"

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
	written, err := store.WriteTo(&buf)
	require.NoError(t, err)

	// Read from snapshot into a new store
	newStore := NewStore[SerializeTestDoc]()
	read, err := newStore.ReadFrom(&buf)
	require.NoError(t, err)

	// Verify byte counts match
	assert.Equal(t, written, read)
	assert.Greater(t, written, int64(0))

	// Verify new store
	newHash, err := newStore.Hash()
	require.NoError(t, err)
	assert.Equal(t, originalHash, newHash)

	retrievedDoc, err := newStore.Get(docs[0].Id)
	require.NoError(t, err)
	assert.Equal(t, docs[0].Data.Name, retrievedDoc.Data.Name)
}
