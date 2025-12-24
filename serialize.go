package docstore

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"
)

const (
	magic = "DSS1"
)

type header struct {
	Magic     string    `json:"magic"`
	Timestamp time.Time `json:"timestamp"`
	DocType   string    `json:"doc_type"`
	Hash      int32     `json:"hash"`
}

// WriteTo writes the store to the given writer.
func (s *Store[T]) WriteTo(w io.Writer) (int64, error) {
	cw := &counter{Writer: w}

	err := func() error {
		hash, err := s.Hash()
		if err != nil {
			return fmt.Errorf("unable to hash documents for snapshot: %w", err)
		}

		typeName := reflect.TypeFor[T]().String()

		encoder := json.NewEncoder(cw)
		encoder.SetIndent("", "  ")

		err = encoder.Encode(header{
			Magic:     magic,
			Timestamp: time.Now().UTC(),
			DocType:   typeName,
			Hash:      hash,
		})
		if err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}

		// Write out each document
		for _, doc := range s.documents {
			if err := encoder.Encode(doc); err != nil {
				return fmt.Errorf("failed to encode document %s: %w", doc.Id, err)
			}
		}
		return nil
	}()

	return cw.written, err
}

// ReadFrom reads the store from the given reader.
func (s *Store[T]) ReadFrom(r io.Reader) (int64, error) {
	cr := &counter{Reader: r}

	typeName := reflect.TypeFor[T]().String()
	decoder := json.NewDecoder(cr)

	// Read out the snapshot info
	header := header{}
	if err := decoder.Decode(&header); err != nil {
		return cr.read, fmt.Errorf("failed to read header: %w", err)
	}

	if header.Magic != magic {
		return cr.read, fmt.Errorf("invalid file format")
	}

	if header.DocType != typeName {
		return cr.read, fmt.Errorf("document type mismatch")
	}

	// Read out each document
	for decoder.More() {
		var doc Document[T]

		if err := decoder.Decode(&doc); err != nil {
			return cr.read, fmt.Errorf("failed to decode document: %w", err)
		}

		s.documents[doc.Id] = doc
	}

	return cr.read, nil
}

type counter struct {
	io.Writer
	io.Reader
	written int64
	read    int64
}

func (c *counter) Write(p []byte) (int, error) {
	n, err := c.Writer.Write(p)
	c.written += int64(n)
	return n, err
}

func (c *counter) Read(p []byte) (int, error) {
	n, err := c.Reader.Read(p)
	c.read += int64(n)
	return n, err
}
