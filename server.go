package docstore

import (
	"encoding/json"
	"net/http"
	"strings"
)

// RegisterHandlers registers the store handlers on a ServeMux.
// This uses the path-based routing available in Go 1.22+.
// For example, if you provide "/api/docs", it will register:
// - POST /api/docs/{id}
// - GET /api/docs/{id}
// - DELETE /api/docs/{id}
func RegisterHandlers[T Document](prefix string, mux *http.ServeMux) {
	prefix = strings.TrimSuffix(prefix, "/")
	mux.Handle("POST "+prefix+"/{id}", PutHandler[T]())
	mux.Handle("GET "+prefix+"/{id}", GetHandler[T]())
	mux.Handle("DELETE "+prefix+"/{id}", DeleteHandler[T]())
}

// PutHandler returns an http.HandlerFunc for adding or updating a document.
// It expects a POST request with a JSON body representing the full document.
// On success, it returns the document (with updated timestamps) and status 200.
func PutHandler[T Document]() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := DocId(r.PathValue("id"))
		if id == EmptyDocId {
			id = GenerateDocId()
		}

		var doc T
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := Put(id, doc); err != nil {
			http.Error(w, "failed to put document: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetHandler returns an http.HandlerFunc for retrieving a document.
// It expects the document ID to be in the URL path (e.g., /docs/my-doc-id).
func GetHandler[T Document]() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := DocId(r.PathValue("id"))
		if id == EmptyDocId {
			http.Error(w, "missing document ID from URL path", http.StatusBadRequest)
			return
		}

		doc, err := GetAs[T](id)
		if err != nil {
			if err == ErrDocumentNotFound {
				http.NotFound(w, r)
			} else if err == ErrDocumentTypeMismatch {
				http.Error(w, "document type mismatch", http.StatusBadRequest)
			} else {
				http.Error(w, "failed to get document: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(*doc)
	}
}

// DeleteHandler returns an http.HandlerFunc for deleting a document.
// It expects the document ID to be in the URL path (e.g., /docs/my-doc-id).
func DeleteHandler[T Document]() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := DocId(r.PathValue("id"))
		if id == EmptyDocId {
			http.Error(w, "missing document ID from URL path", http.StatusBadRequest)
			return
		}

		err := Delete(id)
		if err != nil {
			if err == ErrDocumentNotFound {
				http.NotFound(w, r)
			} else {
				http.Error(w, "failed to delete document: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
