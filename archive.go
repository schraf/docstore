package docstore

import (
	"compress/gzip"
	"errors"
	"os"
)

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
