package docstore

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	// Create an empty store
	store := NewStore[TestDoc]()

	// Create a server for store
	mux := http.NewServeMux()
	docserver := NewServer(store)
	docserver.RegisterHandlers("/testdoc", mux)

	// Create the put doc request
	doc := Document[TestDoc]{
		Id:   GenerateDocId(),
		Data: TestDoc{Name: "John Doe", Age: 30},
	}

	encodedDoc, err := json.Marshal(doc)
	require.NoError(t, err)

	putRequest := httptest.NewRequest(http.MethodPost, "/testdoc", bytes.NewBuffer(encodedDoc))
	putRequest.Header.Set("Content-Type", "application/json")

	// Invoke and validate put request
	putResponse := httptest.NewRecorder()
	mux.ServeHTTP(putResponse, putRequest)
	assert.Equal(t, http.StatusNoContent, putResponse.Code)

	// Create the get request
	getRequest := httptest.NewRequest(http.MethodGet, "/testdoc/"+doc.Id.String(), nil)

	// Invoke and validate get request
	getResponse := httptest.NewRecorder()
	mux.ServeHTTP(getResponse, getRequest)
	assert.Equal(t, http.StatusOK, getResponse.Code)

	var getDoc Document[TestDoc]
	err = json.Unmarshal(getResponse.Body.Bytes(), &getDoc)
	require.NoError(t, err)
	assert.Equal(t, doc.Id, getDoc.Id)
	assert.NotZero(t, getDoc.CreatedAt)
	assert.NotZero(t, getDoc.UpdatedAt)
	assert.Equal(t, doc.Data.Name, getDoc.Data.Name)
	assert.Equal(t, doc.Data.Age, getDoc.Data.Age)

	// Create the delete request
	deleteRequest := httptest.NewRequest(http.MethodDelete, "/testdoc/"+doc.Id.String(), nil)

	// Invoke and validate get request
	deleteResponse := httptest.NewRecorder()
	mux.ServeHTTP(deleteResponse, deleteRequest)
	assert.Equal(t, http.StatusNoContent, deleteResponse.Code)

	// Create the get request after the delete
	getRequest = httptest.NewRequest(http.MethodGet, "/testdoc/"+doc.Id.String(), nil)

	// Invoke and validate get request
	getResponse = httptest.NewRecorder()
	mux.ServeHTTP(getResponse, getRequest)
	assert.Equal(t, http.StatusNotFound, getResponse.Code)
}
