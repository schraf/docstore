# Docstore

`docstore` is a simple, in-memory, key-value document store in Go. It is thread-safe and supports persistence through GOB serialization.

Documents can be of any type. The store is a package-level singleton, providing straightforward functions like `Put`, `Get`, `GetAs`, and `Delete`.

## Features

*   In-memory key-value store.
*   Thread-safe operations.
*   Support for any Go type as a document.
*   Generic-aware `GetAs[T]` for type-safe retrieval.
*   Manual persistence to disk via GOB serialization.
*   GZIP compression for saved stores.
*   A generic HTTP server for exposing a store of a specific type over a REST-like API.

## Getting Started

To use `docstore` in your Go project, you can install it using `go get`:

```bash
go get github.com/schraf/docstore
```

## Usage Example

Because `docstore` uses GOB for serialization, you must register the types of the documents you want to store.

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/schraf/docstore"
)

// Define your document structure.
type MyNote struct {
	Title string
	Body  string
}

func main() {
	// Register your document type for serialization.
	docstore.RegisterType(MyNote{})

	// Create a new document ID.
	id := docstore.GenerateDocId()
	note := MyNote{
		Title: "Hello",
		Body:  "This is my first note in the docstore!",
	}

	// Add the document to the store.
	if err := docstore.Put(id, note); err != nil {
		log.Fatalf("Failed to put document: %v", err)
	}
	fmt.Println("Successfully stored a new note.")

	// Retrieve the document with type assertion.
	retrievedNote, err := docstore.GetAs[MyNote](id)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}
	fmt.Printf("Retrieved note: %+v\n", *retrievedNote)

	// --- Persistence ---

	// Create a temporary file to save the store.
	file, err := os.CreateTemp("", "docstore-*.gz")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// Save the entire store to the file (compressed).
	if err := docstore.WriteAllToFile(file.Name()); err != nil {
		log.Fatalf("Failed to save store: %v", err)
	}
	fmt.Printf("Store saved to %s\n", file.Name())

	// Clear the in-memory store to simulate a restart.
	docstore.Clear()
	fmt.Println("In-memory store cleared.")

	// Load the store back from the file.
	if err := docstore.ReadAllFromFile(file.Name()); err != nil {
		log.Fatalf("Failed to load store: %v", err)
	}
	fmt.Println("Store loaded from file.")

	// Verify the document was restored.
	reloadedNote, err := docstore.GetAs[MyNote](id)
	if err != nil {
		log.Fatalf("Failed to get document after reload: %v", err)
	}
	fmt.Printf("Reloaded note: %+v\n", *reloadedNote)
}
```

### HTTP Server Example

The `Server` is generic and designed to serve a store for a single, specific document type.

```go
package main

import (
	"log"
	"net/http"

	"github.com/schraf/docstore"
)

type UserProfile struct {
	Name  string
	Email string
}

func main() {
	// Register the type for both serialization and the server.
	docstore.RegisterType(UserProfile{})

	// Register the HTTP handlers on a new mux.
	mux := http.NewServeMux()
	RegisterHandlers[UserProfile]("/users", mux)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
```

You can then interact with it using `curl` (requires Go 1.22+ for path-based routing):

```bash
# Create or update a user profile
curl -X POST -d '{"name": "John Doe", "email": "john.doe@example.com"}' http://localhost:8080/users/john

# Retrieve the user profile
curl http://localhost:8080/users/john

# Delete the user profile
curl -X DELETE http://localhost:8080/users/john
```

## Development

This project uses a `Makefile` to automate common development tasks.

### Installing Dependencies

To install the dependencies, run:

```bash
make deps
```

### Running Tests

To run the tests for this project, use the `test` target in the Makefile:

```bash
make test
```

This will run all tests, shuffle their execution order, and run each test 3 times to help detect flaky tests.

### Formatting Code

To format the Go source code, run:

```bash
make fmt
```

### Vetting Code

To run `go vet` on the codebase, use:
```bash
make vet
```
