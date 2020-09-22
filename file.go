package main

import (
	"bufio"
	"os"
	"sync"
)

type dummyReader struct{}

func (d *dummyReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

// BufferedFileReader encapsulate a bufio.Reader with a os.File
type BufferedFileReader struct {
	*bufio.Reader
	file *os.File
}

// Close the file
func (b *BufferedFileReader) Close() error {
	if b.file == nil {
		// Fail silently
		return nil
	}
	return b.file.Close()
}

// Init reinitialize the buffer with the new file
func (b *BufferedFileReader) Init(fileName string) error {
	// If the previous file hasn't been closed, just do it now
	if b.file != nil {
		b.file.Close()
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	b.Reset(file)
	b.file = file
	return nil
}

// NewBufferedFileReader create a new BufferedFileReader
func NewBufferedFileReader(size int) *BufferedFileReader {
	return &BufferedFileReader{
		Reader: bufio.NewReaderSize(&dummyReader{}, size),
	}

}

// NewPoolOfBufferedFileReader create a new sync.Pool of BufferedFileReader
func NewPoolOfBufferedFileReader(size int) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return NewBufferedFileReader(size)
		},
	}
}
