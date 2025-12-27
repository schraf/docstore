package docstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDoc struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestStore_CRUD(t *testing.T) {
	store := NewStore[TestDoc]()

	docId := GenerateDocId()
	doc := Document[TestDoc]{
		Id: docId,
		Data: TestDoc{
			Name: "John Doe",
			Age:  30,
		},
	}

	t.Run("create document", func(t *testing.T) {
		err := store.Put(doc)
		assert.NoError(t, err)
	})

	t.Run("get document", func(t *testing.T) {
		retrieved, err := store.Get(docId)
		assert.NoError(t, err)
		assert.Equal(t, docId, retrieved.Id)
		assert.Equal(t, "John Doe", retrieved.Data.Name)
		assert.Equal(t, 30, retrieved.Data.Age)
	})

	t.Run("get non-existent document", func(t *testing.T) {
		_, err := store.Get(GenerateDocId())
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrDocumentNotFound, err)
	})

	t.Run("update document", func(t *testing.T) {
		updatedDoc := Document[TestDoc]{
			Id: docId,
			Data: TestDoc{
				Name: "John Doe Updated",
				Age:  31,
			},
		}

		err := store.Put(updatedDoc)
		assert.NoError(t, err)

		retrieved, err := store.Get(docId)
		assert.NoError(t, err)
		assert.Equal(t, "John Doe Updated", retrieved.Data.Name)
		assert.Equal(t, 31, retrieved.Data.Age)
	})

	t.Run("delete document", func(t *testing.T) {
		err := store.Delete(docId)
		assert.NoError(t, err)

		_, err = store.Get(docId)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrDocumentNotFound, err)
	})

	t.Run("delete non-existent document", func(t *testing.T) {
		err := store.Delete(GenerateDocId())
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrDocumentNotFound, err)
	})
}

func TestStore_Clear(t *testing.T) {
	store := NewStore[TestDoc]()

	doc := Document[TestDoc]{
		Id:   GenerateDocId(),
		Data: TestDoc{Name: "John Doe", Age: 30},
	}

	err := store.Put(doc)
	assert.NoError(t, err)

	err = store.Clear()
	assert.NoError(t, err)

	getDoc, err := store.Get(doc.Id)
	assert.Nil(t, getDoc)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrDocumentNotFound, err)
}
