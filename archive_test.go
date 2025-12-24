package docstore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertDocumentEqual[T DocData](t *testing.T, expected, actual *Document[T]) {
	t.Helper()
	require.NotNil(t, expected)
	require.NotNil(t, actual)
	assert.Equal(t, expected.Id, actual.Id)
	assert.True(t, expected.CreatedAt.Equal(actual.CreatedAt), "CreatedAt timestamps should be equal")
	assert.True(t, expected.UpdatedAt.Equal(actual.UpdatedAt), "UpdatedAt timestamps should be equal")
	assert.Equal(t, expected.Data, actual.Data)
}

func TestArchive_Roundtrip(t *testing.T) {
	// 1. Create and populate a store
	store := NewStore[TestDoc]()
	doc1 := Document[TestDoc]{
		Id:   GenerateDocId(),
		Data: TestDoc{Name: "John Doe", Age: 30},
	}
	doc2 := Document[TestDoc]{
		Id:   GenerateDocId(),
		Data: TestDoc{Name: "Jane Smith", Age: 25},
	}
	require.NoError(t, store.Put(doc1))
	require.NoError(t, store.Put(doc2))

	// 2. Write the store to an archive in a buffer
	var archiveBuffer bytes.Buffer
	archiveWriter := NewArchiveWriter(&archiveBuffer)
	filename := "test.store"

	err := WriteToArchive(archiveWriter, filename, store)
	require.NoError(t, err)
	require.NoError(t, archiveWriter.Close())

	// 3. Read the store back from the archive
	archiveReader, err := NewArchiveReader(&archiveBuffer)
	require.NoError(t, err)
	defer archiveReader.Close()

	readStore, err := ReadFromArchive[TestDoc](archiveReader, filename)
	require.NoError(t, err)
	require.NotNil(t, readStore)

	// 4. Verify the contents of the restored store
	originalResult, err := store.Select(Query[TestDoc]{})
	require.NoError(t, err)
	readResult, err := readStore.Select(Query[TestDoc]{})
	require.NoError(t, err)
	assert.Equal(t, originalResult.Total, readResult.Total)

	originalDoc1, err := store.Get(doc1.Id)
	require.NoError(t, err)
	retrievedDoc1, err := readStore.Get(doc1.Id)
	require.NoError(t, err)
	assertDocumentEqual(t, originalDoc1, retrievedDoc1)

	originalDoc2, err := store.Get(doc2.Id)
	require.NoError(t, err)
	retrievedDoc2, err := readStore.Get(doc2.Id)
	require.NoError(t, err)
	assertDocumentEqual(t, originalDoc2, retrievedDoc2)
}

func TestArchive_FileNotFound(t *testing.T) {
	// Create an empty archive
	var archiveBuffer bytes.Buffer
	archiveWriter := NewArchiveWriter(&archiveBuffer)
	store := NewStore[TestDoc]()
	err := WriteToArchive(archiveWriter, "some-other-file.store", store)
	require.NoError(t, err)
	require.NoError(t, archiveWriter.Close())

	// Attempt to read a file that doesn't exist
	archiveReader, err := NewArchiveReader(&archiveBuffer)
	require.NoError(t, err)
	defer archiveReader.Close()

	_, err = ReadFromArchive[TestDoc](archiveReader, "non-existent-file.store")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), `file "non-existent-file.store" not found in archive`)
}
