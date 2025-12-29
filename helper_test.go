package docstore

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestDoc struct {
	Name string
	Age  int
}

func AssertDoc(t *testing.T, id DocId, expected TestDoc) {
	t.Helper()

	doc, err := GetAs[TestDoc](id)
	require.NoError(t, err)
	assert.Equal(t, expected.Name, doc.Name)
	assert.Equal(t, expected.Age, doc.Age)
}

func AssertNoDoc(t *testing.T, id DocId) {
	t.Helper()

	_, err := Get(id)
	assert.ErrorIs(t, ErrDocumentNotFound, err)
}

func AddTestDoc(t *testing.T, name string, age int) (DocId, TestDoc) {
	t.Helper()

	id := GenerateDocId()
	doc := TestDoc{
		Name: name,
		Age:  age,
	}

	err := Put(id, doc)
	require.NoError(t, err)

	return id, doc
}

func RequireDocPutRequest(t *testing.T, name string, age int) (DocId, TestDoc, *http.Request) {
	t.Helper()

	id := GenerateDocId()
	doc := TestDoc{
		Name: name,
		Age:  age,
	}

	encodedDoc, err := json.Marshal(doc)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/testdoc/"+id.String(), bytes.NewBuffer(encodedDoc))
	request.Header.Set("Content-Type", "application/json")

	return id, doc, request
}

func RequireDocBulkPutRequest(t *testing.T, docs ...TestDoc) ([]DocId, *http.Request) {
	t.Helper()

	bulk := map[string]TestDoc{}
	docIds := []DocId{}

	for _, doc := range docs {
		id := GenerateDocId()
		bulk[id.String()] = doc
		docIds = append(docIds, id)
	}

	encodedDocs, err := json.Marshal(bulk)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/testdoc/bulk", bytes.NewBuffer(encodedDocs))
	request.Header.Set("Content-Type", "application/json")

	return docIds, request
}

func AssertResponseCode(t *testing.T, mux *http.ServeMux, request *http.Request, code int) {
	t.Helper()

	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	assert.Equal(t, code, response.Code)
}

func RequireResponseBody[T Document](t *testing.T, mux *http.ServeMux, request *http.Request) T {
	t.Helper()

	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	require.Equal(t, http.StatusOK, response.Code)

	var doc T
	err := json.Unmarshal(response.Body.Bytes(), &doc)
	require.NoError(t, err)

	return doc
}

func TempFilename() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), "docstore_"+hex.EncodeToString(randBytes))
}
