package docstore

import (
	"compress/gzip"
	"errors"
	"os"
)

// WriteAllToFile writes the entire docstore to a single file.
func WriteAllToFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	gzipWriter := gzip.NewWriter(file)

	defer func() {
		err = errors.Join(err, gzipWriter.Close(), file.Close())
	}()

	err = WriteAll(gzipWriter)

	return
}

// ReadAllFromFile reads the entire docstore from a single file, replacing any
// contents that were in the docstore.
func ReadAllFromFile(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		file.Close()
		return err
	}

	defer func() {
		err = errors.Join(err, gzipReader.Close(), file.Close())
	}()

	err = ReadAll(gzipReader)

	return
}
