package docstore

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	Clear()

	// setup the test server
	mux := http.NewServeMux()
	docserver := NewServer[TestDoc]()
	docserver.RegisterHandlers("/testdoc", mux)

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
