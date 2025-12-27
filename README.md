# docstore

A simple in-memory document store written in Go.

## Description

`docstore` is a lightweight, generic, in-memory document storage library. It allows you to store and retrieve documents in a structured way. The library is designed to be simple to use and integrate into any Go project.

## Features

- **Generic:** Can store any type of document data.
- **In-Memory:** Fast and efficient for small to medium-sized datasets.
- **CRUD Operations:** Supports Create, Read, Update, and Delete operations.
- **Serialization:** Save and load the document store to and from a file.
- **Archiving:** Utilities to bundle a store into a gzipped tarball.

## Components

- **`Store`:** The main component that manages the documents.
- **`Document`:** A wrapper around your data that includes metadata like ID, and timestamps.
- **`ArchiveWriter`:** Used to serialize multiple Store components into a archive file.
- **`ArchiveReader`:** Used to deserialize Stores from an archive file that was written using `ArchiveWriter`.

## Installation

```bash
go get github.com/schraf/docstore
```

## Usage

Here is a simple example of how to use `docstore`:

```go
package main

import (
	"fmt"
	"log"

	"github.com/schraf/docstore"
)

// Define your document data structure
type MyDoc struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// Create a new store
	store := docstore.NewStore[MyDoc]()

	// Create a new document
	doc := docstore.Document[MyDoc]{
		Id: docstore.GenerateDocId(),
		Data: MyDoc{
			Name: "John Doe",
			Age:  30,
		},
	}

	// Add the document to the store
	if err := store.Put(doc); err != nil {
		log.Fatalf("Failed to put document: %v", err)
	}

	// Retrieve the document
	retrievedDoc, err := store.Get(doc.Id)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}

	fmt.Printf("Retrieved document: %+v\n", retrievedDoc)
}
```

## Development

This project uses a `Makefile` for common development tasks:

-   `make test`: Runs all unit tests.
-   `make vet`: Vets the code for suspicious constructs.
-   `make fmt`: Formats the code according to Go standards.
-   `make all`: Runs `vet` and `test`.
-   `make deps`: Installs dependencies.
-   `make help`: Shows all available commands.

To run tests, simply execute:

```bash
make test
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

