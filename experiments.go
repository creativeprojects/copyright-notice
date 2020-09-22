package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

const (
	bufferDefaultSize = 16384
)

var (
	bufReaderPool *sync.Pool
	bufferPool    *sync.Pool
)

func init() {
	bufReaderPool = &sync.Pool{
		New: func() interface{} {
			return newSourceFileReader(
				// Creates a reader using stdin. Actually stdin will never be read but a valid default is needed
				bufio.NewReaderSize(os.Stdin, bufferDefaultSize),
				ioutil.NopCloser(os.Stdin),
			)
		},
	}
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
}

type sourceFileReader struct {
	reader *bufio.Reader
	closer io.Closer
}

func (s *sourceFileReader) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}
func (s *sourceFileReader) Close() error {
	return s.closer.Close()
}
func (s *sourceFileReader) InitSeparate(reader io.Reader, closer io.Closer) {
	s.reader.Reset(reader)
	s.closer = closer
}
func (s *sourceFileReader) Init(file *os.File) {
	// Turns out the Reset method returns a new instance of bufio.Reader (but keep the inside []byte)
	s.reader.Reset(file)
	s.closer = file
}

func newSourceFileReader(reader *bufio.Reader, closer io.Closer) io.ReadCloser {
	return &sourceFileReader{
		reader: reader,
		closer: closer,
	}
}

func getFileReader(fileName string) (io.ReadCloser, error) {
	file, err := os.Open(fileName)
	if err != nil {
		progress(fileName, fileStatusCannotOpen, err)
		return nil, err
	}

	buffer := bufio.NewReaderSize(file, bufferDefaultSize)
	bom, err := buffer.Peek(3)
	if err != nil {
		progress(fileName, fileStatusError, err)
		return nil, err
	}
	if bom[0] == 0xef && bom[1] == 0xbb && bom[2] == 0xbf {
		// This is a bom, move the file forward 3 positions
		_, err := buffer.Discard(3)
		if err != nil {
			progress(fileName, fileStatusError, err)
			return nil, err
		}
	}
	return newSourceFileReader(buffer, file), nil
}

func getFileReaderFromPool(fileName string) (io.ReadCloser, error) {
	file, err := os.Open(fileName)
	if err != nil {
		progress(fileName, fileStatusCannotOpen, err)
		return nil, err
	}

	fileReader := bufReaderPool.Get().(*sourceFileReader)
	fileReader.Init(file)
	bom, err := fileReader.reader.Peek(3)
	if err != nil {
		progress(fileName, fileStatusError, err)
		return nil, err
	}
	if bom[0] == 0xef && bom[1] == 0xbb && bom[2] == 0xbf {
		// This is a bom, move the file forward 3 positions
		_, err := fileReader.reader.Discard(3)
		if err != nil {
			progress(fileName, fileStatusError, err)
			return nil, err
		}
	}
	return fileReader, nil
}

func readFileIntoBuffer(fileName string, buffer *bytes.Buffer) error {
	var err error
	if buffer == nil {
		panic("buffer parameter cannot be nil")
	}

	file, err := os.Open(fileName)
	if err != nil {
		progress(fileName, fileStatusCannotOpen, err)
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		progress(fileName, fileStatusCannotOpen, err)
		return err
	}
	if info.Size() > maxFileSize {
		progress(fileName, fileStatusTooBig, nil)
		return err
	}
	fileSize := info.Size()

	// Read the first 3 bytes and see if it's a BOM.
	// If it is, we're going to ignore it (and save the file without)
	// Don't worry, Visual Studio will put it back eventually
	// var bom [3]byte
	// _, err = io.ReadFull(file, bom[:])
	// if err != nil {
	// 	progress(fileName, fileStatusError, err)
	// 	return err
	// }
	// if bom[0] != 0xef || bom[1] != 0xbb || bom[2] != 0xbf {
	// 	_, err = file.Seek(0, 0) // Not a BOM -- seek back to the beginning
	// 	if err != nil {
	// 		progress(fileName, fileStatusError, err)
	// 		return err
	// 	}
	// } else {
	// 	// Forget the first 3 bytes
	// 	fileSize -= 3
	// }

	// buffer := &bytes.Buffer{}
	bytesRead, err := buffer.ReadFrom(file)
	if err != nil && err != io.EOF {
		progress(fileName, fileStatusError, err)
		return err
	}
	if err == io.EOF || bytesRead != fileSize {
		progress(fileName, fileStatusError, fmt.Errorf("file size = %d bytes but read %d bytes instead", fileSize, bytesRead))
		return err
	}
	return nil
}
