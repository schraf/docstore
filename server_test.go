package docstore

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	Clear()

	// setup the test server
	mux := http.NewServeMux()
	RegisterHandlers[TestDoc]("/testdoc", mux)

	// add a document
	jimId, jim, request := RequireDocPutRequest(t, "jim", 22)
	AssertResponseCode(t, mux, request, http.StatusNoContent)
	AssertDoc(t, jimId, jim)

	// get a document
	request = httptest.NewRequest(http.MethodGet, "/testdoc/"+jimId.String(), nil)
	responseDoc := RequireResponseBody[TestDoc](t, mux, request)
	assert.Equal(t, jim.Name, responseDoc.Name)
	assert.Equal(t, jim.Age, responseDoc.Age)

	// delete a document
	request = httptest.NewRequest(http.MethodDelete, "/testdoc/"+jimId.String(), nil)
	AssertResponseCode(t, mux, request, http.StatusNoContent)
	AssertNoDoc(t, jimId)

	// validate no doc
	request = httptest.NewRequest(http.MethodGet, "/testdoc/"+jimId.String(), nil)
	AssertResponseCode(t, mux, request, http.StatusNotFound)
}

func TestServer_List(t *testing.T) {
	Clear()

	jimId, _ := AddTestDoc(t, "jim", 22)
	joeId, _ := AddTestDoc(t, "joe", 32)
	bobId, _ := AddTestDoc(t, "bob", 42)

	// setup the test server
	mux := http.NewServeMux()
	RegisterHandlers[TestDoc]("/testdoc", mux)

	request := httptest.NewRequest(http.MethodGet, "/testdoc", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	var docs []struct {
		Id   string  `json:"id"`
		Data TestDoc `json:"data"`
	}

	err := json.Unmarshal(response.Body.Bytes(), &docs)
	require.NoError(t, err)
	assert.Equal(t, 3, len(docs))

	jimCount := 0
	joeCount := 0
	bobCount := 0

	for _, doc := range docs {
		switch doc.Id {
		case jimId.String():
			jimCount++
		case joeId.String():
			joeCount++
		case bobId.String():
			bobCount++
		default:
			t.Fatalf("unexpected doc '%s'", doc.Id)
		}
	}

	assert.Equal(t, 1, jimCount)
	assert.Equal(t, 1, joeCount)
	assert.Equal(t, 1, bobCount)
}
