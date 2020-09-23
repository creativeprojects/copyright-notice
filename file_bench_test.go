package main

import (
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

const (
	bufferSize = 16 * 1024
)

//
// Results
// ========
//
// There's barely any difference in bettween these two benchmarks:
// the time it takes to allocate memory is nothing compared to the time it takes to read from the disk
//

func BenchmarkFileWithBufferedFileReader(b *testing.B) {
	var err error
	fileName := ""
	switch runtime.GOOS {
	case "windows":
		fileName = `C:\WINDOWS\notepad.exe`
	case "darwin":
		fileName = `/System/Library/Fonts/Helvetica.ttc`
	}
	if fileName == "" {
		b.Skip("file not found")
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		b.Skip(err)
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := NewBufferedFileReader(bufferSize)
		err = reader.Init(fileName)
		if err != nil {
			b.Skip(err)
		}
		copied, err := io.Copy(ioutil.Discard, reader)
		if err != nil {
			b.Error(err)
		}
		if copied != fileInfo.Size() {
			b.Errorf("file size %d, but read %d bytes", fileInfo.Size(), copied)
		}
		reader.Close()
	}
}

func BenchmarkFileWithPoolOfBufferedFileReader(b *testing.B) {
	var err error
	fileName := ""
	switch runtime.GOOS {
	case "windows":
		fileName = `C:\WINDOWS\notepad.exe`
	case "darwin":
		fileName = `/System/Library/Fonts/Helvetica.ttc`
	}
	if fileName == "" {
		b.Skip("file not found")
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		b.Skip(err)
	}

	b.ReportAllocs()

	pool := NewPoolOfBufferedFileReader(bufferSize)

	for i := 0; i < b.N; i++ {
		reader := pool.Get().(*BufferedFileReader)
		err = reader.Init(fileName)
		if err != nil {
			b.Skip(err)
		}
		copied, err := io.Copy(ioutil.Discard, reader)
		if err != nil {
			b.Error(err)
		}
		if copied != fileInfo.Size() {
			b.Errorf("file size %d, but read %d bytes", fileInfo.Size(), copied)
		}
		reader.Close()
		pool.Put(reader)
	}
}
