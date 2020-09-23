package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	name    string
	size    int
	content []byte
	ready   bool
}

func NewFile(bufferSize int) *File {
	return &File{
		content: make([]byte, 0, bufferSize),
	}
}

func (f *File) Reset() *File {
	f.name = ""
	f.size = 0
	f.content = f.content[:0]
	f.ready = false
	return f
}

func (f *File) Read(name string, size int64) error {
	const maxInt = 2147483647
	if f.ready {
		// clear up the buffer first
		f.Reset()
	}
	if name == "" {
		return NewError(FileErrorInvalidName, nil)
	}
	// max value for an int
	if f.size > maxInt {
		return NewError(FileErrorTooBig, errors.New("file size (int64) doesn't fit in an int32"))
	}

	f.name = name
	f.size = int(size)

	if size == 0 {
		// file is empty, nothing to read
		f.ready = true
		return nil
	}
	if f.size > cap(f.content) {
		return NewError(FileErrorTooBig, fmt.Errorf("file size = %d bigger than buffer size = %d", f.size, cap(f.content)))
	}
	file, err := os.Open(f.name)
	if err != nil {
		return NewError(FileErrorCannotOpen, err)
	}
	defer file.Close()

	// reslice the buffer
	f.content = f.content[:f.size]
	// and read the whole file
	read, err := file.Read(f.content)
	if err != nil && err != io.EOF {
		return NewError(FileErrorReading, err)
	}
	if err == io.EOF || read != f.size {
		return NewError(FileErrorReading, fmt.Errorf("file size = %d bytes but read %d bytes instead", f.size, read))
	}

	f.ready = true
	return nil
}

func (f *File) IsReady() bool {
	return f.ready
}

func (f *File) HasUTF8BOM() bool {
	return hasUTF8BOM(f.content)
}

// Bytes returns the file content (with the UTF8 BOM stripped out if any)
func (f *File) Bytes() []byte {
	if f.HasUTF8BOM() {
		return f.content[3:]
	}
	return f.content
}

// AddHeader saves the file with the new header.
// Instead of creating a file in place, it saves a temporary file then renames it
func (f *File) AddHeader(header []byte, keepUTF8BOM bool) error {
	var err error
	randomBytes := make([]byte, 10)
	randomGenerator.Read(randomBytes)
	tempFilename := filepath.Join(filepath.Dir(f.name), "$"+fmt.Sprintf("%x", randomBytes)+"$"+filepath.Base(f.name))

	err = f.saveFile(tempFilename, header, keepUTF8BOM)
	if err != nil {
		return err
	}
	// Move the temp file into place
	err = os.Rename(tempFilename, f.name)
	if err != nil {
		// Try to delete the temp file
		os.Remove(tempFilename)
		return err
	}
	return nil
}

func (f *File) saveFile(filename string, header []byte, keepUTF8BOM bool) error {
	outputFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return f.saveContent(outputFile, header, keepUTF8BOM)
}

func (f *File) saveContent(writer io.Writer, header []byte, keepUTF8BOM bool) error {
	var err error

	// Write the BOM if it was present
	if keepUTF8BOM && f.HasUTF8BOM() {
		_, err := writer.Write(UTF8BOM)
		if err != nil {
			return err
		}
	}

	// Write the copyright notice
	_, err = writer.Write(header)
	if err != nil {
		return err
	}

	// Then write the file content
	writer.Write(f.Bytes())
	if err != nil {
		return err
	}
	return nil
}
