package docstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash_Consistent(t *testing.T) {
	store := NewStore[TestDoc]()

	doc1 := Document[TestDoc]{Id: NewDocId("1"), Data: TestDoc{Name: "Alice", Age: 30}}
	doc2 := Document[TestDoc]{Id: NewDocId("2"), Data: TestDoc{Name: "Bob", Age: 25}}

	err := store.Put(doc1)
	assert.NoError(t, err)
	err = store.Put(doc2)
	assert.NoError(t, err)

	hash1, err := store.Hash()
	assert.NoError(t, err)

	hash2, err := store.Hash()
	assert.NoError(t, err)

	assert.Equal(t, hash1, hash2, "Hash should be consistent for the same store state")
}
