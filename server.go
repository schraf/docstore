package docstore

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Server is an HTTP server for a Store.
type Server[T DocData] struct {
	store *Store[T]
}

// NewServer creates a new server for a store.
func NewServer[T DocData](store *Store[T]) *Server[T] {
	return &Server[T]{
		store: store,
	}
}

// RegisterHandlers registers the store handlers on a ServeMux.
// This uses the path-based routing available in Go 1.22+.
// For example, if you provide "/api/docs", it will register:
// - POST /api/docs
// - GET /api/docs/{id}
// - DELETE /api/docs/{id}
func (s *Server[T]) RegisterHandlers(prefix string, mux *http.ServeMux) {
	prefix = strings.TrimSuffix(prefix, "/")
	mux.Handle("POST "+prefix, s.PutHandler())
	mux.Handle("GET "+prefix+"/{id}", s.GetHandler())
	mux.Handle("DELETE "+prefix+"/{id}", s.DeleteHandler())
}

// PutHandler returns an http.HandlerFunc for adding or updating a document.
// It expects a POST request with a JSON body representing the full document.
// On success, it returns the document (with updated timestamps) and status 200.
func (s *Server[T]) PutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var doc Document[T]
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		if doc.Id == EmptyDocId {
			doc.Id = GenerateDocId()
		}

		if err := s.store.Put(doc); err != nil {
			http.Error(w, "failed to put document: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetHandler returns an http.HandlerFunc for retrieving a document.
// It expects the document ID to be in the URL path (e.g., /docs/my-doc-id).
func (s *Server[T]) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := DocId(r.PathValue("id"))
		if id == EmptyDocId {
			http.Error(w, "missing document ID from URL path", http.StatusBadRequest)
			return
		}

		doc, err := s.store.Get(id)
		if err != nil {
			if err == ErrDocumentNotFound {
				http.NotFound(w, r)
			} else {
				http.Error(w, "failed to get document: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(doc)
	}
}

// DeleteHandler returns an http.HandlerFunc for deleting a document.
// It expects the document ID to be in the URL path (e.g., /docs/my-doc-id).
func (s *Server[T]) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := DocId(r.PathValue("id"))
		if id == EmptyDocId {
			http.Error(w, "missing document ID from URL path", http.StatusBadRequest)
			return
		}

		err := s.store.Delete(id)
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
