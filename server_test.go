package docstore

import (
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

func TestServer_BulkPut(t *testing.T) {
	Clear()

	jim := TestDoc{Name: "jim", Age: 22}
	joe := TestDoc{Name: "joe", Age: 32}
	bob := TestDoc{Name: "bob", Age: 42}

	// setup the test server
	mux := http.NewServeMux()
	RegisterHandlers[TestDoc]("/testdoc", mux)

	// create the bulk request
	ids, request := RequireDocBulkPutRequest(t, jim, joe, bob)
	AssertResponseCode(t, mux, request, http.StatusNoContent)
	require.Equal(t, 3, len(ids))

	AssertDoc(t, ids[0], jim)
	AssertDoc(t, ids[1], joe)
	AssertDoc(t, ids[2], bob)
}
