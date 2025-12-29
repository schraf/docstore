package docstore

import (
	"encoding/gob"
	"fmt"
	"io"
)

const (
	magic = "DSS1"
)

type header struct {
	Magic string
	Count int
}

type document struct {
	Id  DocId
	Doc Document
}

func RegisterType(d Document) {
	gob.Register(d)
}

// WriteAll writes all of the documents to the given writer.
func WriteAll(w io.Writer) error {
	encoder := gob.NewEncoder(w)

	storeLock.RLock()
	defer storeLock.RUnlock()

	err := encoder.Encode(header{
		Magic: magic,
		Count: len(store),
	})
	if err != nil {
		return err
	}

	for id, doc := range store {
		if err := encoder.Encode(document{
			Id:  id,
			Doc: doc,
		}); err != nil {
			return err
		}
	}

	return nil
}

// ReadAll reads all of the documents from the given reader.
func ReadAll(r io.Reader) error {
	decoder := gob.NewDecoder(r)

	// Read out the snapshot info
	var h header
	if err := decoder.Decode(&h); err != nil {
		return err
	}

	if h.Magic != magic || h.Count < 0 {
		return fmt.Errorf("invalid file header")
	}

	// Read out each document
	for i := 0; i < h.Count; i++ {
		var doc document

		if err := decoder.Decode(&doc); err != nil {
			return err
		}

		Put(doc.Id, doc.Doc)
	}

	return nil
}
