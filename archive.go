package docstore

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"time"
)

// ArchiveWriter is a writer for creating tar.gz archives.
type ArchiveWriter struct {
	archive *tar.Writer
	closer  func() error
}

// NewArchiveWriter creates a new archive writer.
func NewArchiveWriter(w io.Writer) *ArchiveWriter {
	gzipWriter := gzip.NewWriter(w)
	tarWriter := tar.NewWriter(gzipWriter)

	return &ArchiveWriter{
		archive: tarWriter,
		closer: func() error {
			return errors.Join(tarWriter.Close(), gzipWriter.Close())
		},
	}
}

// Close closes the archive writer.
func (a *ArchiveWriter) Close() error {
	return a.closer()
}

// WriteToArchive writes a store to the archive.
func WriteToArchive[T DocData](w *ArchiveWriter, filename string, store *Store[T]) error {
	var buffer bytes.Buffer

	_, err := store.WriteTo(&buffer)
	if err != nil {
		return fmt.Errorf("failed writing store: %w", err)
	}

	data, err := io.ReadAll(&buffer)
	if err != nil {
		return fmt.Errorf("failed reading serialized store: %w", err)
	}

	header := tar.Header{
		Name:    filename,
		Mode:    0666,
		Size:    int64(len(data)),
		ModTime: time.Now(),
	}

	if err := w.archive.WriteHeader(&header); err != nil {
		return err
	}

	if _, err := w.archive.Write(data); err != nil {
		return err
	}

	return nil
}

// ArchiveReader is a reader for reading tar.gz archives.
type ArchiveReader struct {
	archive *tar.Reader
	closer  func() error
}

// NewArchiveReader creates a new archive reader.
func NewArchiveReader(r io.Reader) (*ArchiveReader, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed creating gzip reader: %w", err)
	}

	tarReader := tar.NewReader(gzipReader)

	return &ArchiveReader{
		archive: tarReader,
		closer:  gzipReader.Close,
	}, nil
}

// Close closes the archive reader.
func (a *ArchiveReader) Close() error {
	return a.closer()
}

// ReadFromArchive reads a store from the archive.
func ReadFromArchive[T DocData](r *ArchiveReader, filename string) (*Store[T], error) {
	for {
		header, err := r.archive.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read archive header: %w", err)
		}

		if header.Name == filename {
			store := NewStore[T]()
			if _, err := store.ReadFrom(r.archive); err != nil {
				return nil, fmt.Errorf("failed to read store from archive file %q: %w", filename, err)
			}
			return store, nil
		}
	}

	return nil, fmt.Errorf("file %q not found in archive", filename)
}
