package docstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestStore_Select(t *testing.T) {
	store := NewStore[TestDoc]()

	// Create test users
	docs := []Document[TestDoc]{
		{Id: GenerateDocId(), Data: TestDoc{Name: "John Doe", Age: 30}},
		{Id: GenerateDocId(), Data: TestDoc{Name: "Jane Smith", Age: 25}},
		{Id: GenerateDocId(), Data: TestDoc{Name: "Bob Johnson", Age: 35}},
	}

	for _, doc := range docs {
		err := store.Put(doc)
		require.NoError(t, err)
	}

	t.Run("query all", func(t *testing.T) {
		result, err := store.Select(Query[TestDoc]{})
		assert.NoError(t, err)
		assert.Len(t, result.Documents, 3)
		assert.Equal(t, 3, result.Total)
	})

	t.Run("query with filters", func(t *testing.T) {
		ageFilter := func(d *Document[TestDoc]) bool { return d.Data.Age > 28 }
		result, err := store.Select(Query[TestDoc]{
			Filters: []QueryFilter[TestDoc]{ageFilter},
		})
		assert.NoError(t, err)
		assert.Len(t, result.Documents, 2)
		assert.Equal(t, 2, result.Total)
	})

	t.Run("query with limit", func(t *testing.T) {
		result, err := store.Select(Query[TestDoc]{Limit: 2})
		assert.NoError(t, err)
		assert.Len(t, result.Documents, 2)
		assert.Equal(t, 3, result.Total)
	})

	t.Run("query with sortby", func(t *testing.T) {
		// Sort by age ascending
		ageSortAsc := QuerySort[TestDoc](func(a, b *Document[TestDoc]) bool { return a.Data.Age < b.Data.Age })
		result, err := store.Select(Query[TestDoc]{
			SortBy: &ageSortAsc,
		})
		assert.NoError(t, err)
		require.Len(t, result.Documents, 3)
		assert.Equal(t, "Jane Smith", result.Documents[0].Data.Name)
		assert.Equal(t, "John Doe", result.Documents[1].Data.Name)
		assert.Equal(t, "Bob Johnson", result.Documents[2].Data.Name)

		// Sort by age descending
		ageSortDesc := QuerySort[TestDoc](func(a, b *Document[TestDoc]) bool { return a.Data.Age > b.Data.Age })
		result, err = store.Select(Query[TestDoc]{
			SortBy: &ageSortDesc,
		})
		assert.NoError(t, err)
		require.Len(t, result.Documents, 3)
		assert.Equal(t, "Bob Johnson", result.Documents[0].Data.Name)
		assert.Equal(t, "John Doe", result.Documents[1].Data.Name)
		assert.Equal(t, "Jane Smith", result.Documents[2].Data.Name)
	})

	t.Run("query with filter and limit", func(t *testing.T) {
		ageFilter := func(d *Document[TestDoc]) bool { return d.Data.Age > 28 } // John Doe (30), Bob Johnson (35)
		result, err := store.Select(Query[TestDoc]{
			Filters: []QueryFilter[TestDoc]{ageFilter},
			Limit:   1,
		})
		assert.NoError(t, err)
		require.Len(t, result.Documents, 1)
		assert.Equal(t, 2, result.Total) // Total should still be 2 before limiting
		assert.NotEqual(t, "Jane Smith", result.Documents[0].Data.Name)
	})

	t.Run("query with sort and limit", func(t *testing.T) {
		ageSortAsc := QuerySort[TestDoc](func(a, b *Document[TestDoc]) bool { return a.Data.Age < b.Data.Age })
		result, err := store.Select(Query[TestDoc]{
			SortBy: &ageSortAsc,
			Limit:  2,
		})
		assert.NoError(t, err)
		require.Len(t, result.Documents, 2)
		assert.Equal(t, 3, result.Total) // Total should still be 3 before limiting
		assert.Equal(t, "Jane Smith", result.Documents[0].Data.Name)
		assert.Equal(t, "John Doe", result.Documents[1].Data.Name)
	})

	t.Run("query with filter and sort", func(t *testing.T) {
		ageFilter := func(d *Document[TestDoc]) bool { return d.Data.Age > 28 } // John Doe (30), Bob Johnson (35)
		ageSortAsc := QuerySort[TestDoc](func(a, b *Document[TestDoc]) bool { return a.Data.Age < b.Data.Age })
		result, err := store.Select(Query[TestDoc]{
			Filters: []QueryFilter[TestDoc]{ageFilter},
			SortBy:  &ageSortAsc,
		})
		assert.NoError(t, err)
		assert.Len(t, result.Documents, 2)
		require.Equal(t, 2, result.Total)
		assert.Equal(t, "John Doe", result.Documents[0].Data.Name)
		assert.Equal(t, "Bob Johnson", result.Documents[1].Data.Name)
	})
}
