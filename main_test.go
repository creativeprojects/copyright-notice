package main

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/creativeprojects/clog"
)

func init() {
	flags.extensions = []string{".go"}
	cleanupConfiguration()

	excludeList := []string{
		"**/.*",
		"**/$*",
		"bin",
		"vendor",
		"packages",
		"node_modules",
	}
	exclusions := newExclusion(excludeList...)

	parseDir := "."
	// gopath := build.Default.GOPATH
	// if gopath != "" {
	// 	parseDir = gopath
	// }
	clog.SetDefaultLogger(clog.NewLogger(clog.NewDiscardHandler()))
	parseDirectory(parseDir, exclusions, func(int) {}, func() {})
}

func TestMaxSize(t *testing.T) {
	t.Log("Found", fileQueue.Len(), "files")
	t.Log("Max file size", maxSize, "kb")
	t.Log("Average file size", int(maxSize)/fileQueue.Len(), "kb")
}

// New Results with isolated runs:
// BenchmarkReadIntoByteSlice-12                       7218	    162098 ns/op	   11413 B/op	       6 allocs/op
// BenchmarkReadIntoBytesBuffer-12                     7282	    148917 ns/op	   34374 B/op	       6 allocs/op
// BenchmarkReadIntoBytesBufferWithPool-12             7011	    147481 ns/op	     553 B/op	       4 allocs/op
// BenchmarkReadFromBufIOToBuffer-12                   6920	    151152 ns/op	   48415 B/op	       9 allocs/op
// BenchmarkReadFromBufIOToPoolOfBuffer-12             7004	    146264 ns/op	   17004 B/op	       6 allocs/op
// BenchmarkReadFromBufIOFromPoolToBuffer-12           7347	    151984 ns/op	   34506 B/op	       6 allocs/op
// BenchmarkReadFromBufIOFromPoolToPoolOfBuffer-12     7639	    153574 ns/op	     629 B/op	       3 allocs/op
// BenchmarkReadFromBufIOToByteSlice-12                7045	    149583 ns/op	   28973 B/op	      18 allocs/op
// BenchmarkReadFromBufIOFromPoolToByteSlice-12        7575	    144371 ns/op	   14125 B/op	      15 allocs/op

// Results
// ========
// BenchmarkReadIntoByteSlice-6                               24384             47331 ns/op            8069 B/op          9 allocs/op
// BenchmarkReadIntoBytesBuffer-6                             21230             54550 ns/op           19990 B/op         10 allocs/op
// BenchmarkReadIntoBytesBufferWithPool-6                     26488             46162 ns/op            1555 B/op          7 allocs/op
// BenchmarkReadIntoBufIO-6                                   30105             40449 ns/op           17918 B/op          9 allocs/op
// BenchmarkReadIntoBufIOFromPool-6                           35013             33851 ns/op            1382 B/op          6 allocs/op

func BenchmarkReadIntoByteSlice(b *testing.B) {
	b.ReportAllocs()
	if fileQueue.Len() == 0 {
		b.Skip("No source file")
	}
	e := fileQueue.Front()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileEntry := e.Value.(FileEntry)
		buffer, err := readFile(fileEntry.Name)
		if err != nil {
			b.Log("Cannot read file", fileEntry.Name, err)
		} else if len(buffer) == 0 {
			b.Log("Empty file", fileEntry.Name)
		}
		e = e.Next()
		if e == nil {
			e = fileQueue.Front()
		}
	}
}

func BenchmarkReadIntoBytesBuffer(b *testing.B) {
	b.ReportAllocs()
	if fileQueue.Len() == 0 {
		b.Skip("No source file")
	}
	e := fileQueue.Front()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileEntry := e.Value.(FileEntry)
		buffer := &bytes.Buffer{}
		if buffer.Cap() < int(fileEntry.Size) {
			buffer.Grow(int(fileEntry.Size) - buffer.Cap() + 1)
		}
		err := readFileIntoBuffer(fileEntry.Name, buffer)
		if err != nil {
			b.Log("Cannot read file", fileEntry.Name, err)
		} else if buffer.Len() == 0 {
			b.Log("Empty file", fileEntry.Name)
		}
		e = e.Next()
		if e == nil {
			e = fileQueue.Front()
		}
	}
}

func BenchmarkReadIntoBytesBufferWithPool(b *testing.B) {
	b.ReportAllocs()
	if fileQueue.Len() == 0 {
		b.Skip("No source file")
	}
	localBufferPool := sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	e := fileQueue.Front()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileEntry := e.Value.(FileEntry)
		buffer := localBufferPool.Get().(*bytes.Buffer)
		buffer.Reset()
		if buffer.Cap() < int(fileEntry.Size) {
			buffer.Grow(int(fileEntry.Size) - buffer.Cap() + 1)
		}
		err := readFileIntoBuffer(fileEntry.Name, buffer)
		if err != nil {
			b.Log("Cannot read file", fileEntry.Name, err)
		} else if buffer.Len() == 0 {
			b.Log("Empty file", fileEntry.Name)
		}
		localBufferPool.Put(buffer)
		e = e.Next()
		if e == nil {
			e = fileQueue.Front()
		}
	}
}

func BenchmarkReadFromBufIOToBuffer(b *testing.B) {
	b.ReportAllocs()
	if fileQueue.Len() == 0 {
		b.Skip("No source file")
	}
	e := fileQueue.Front()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileEntry := e.Value.(FileEntry)
		reader, err := getFileReader(fileEntry.Name)
		if err != nil {
			b.Log("Cannot open file", fileEntry.Name, err)
		} else {
			buffer := &bytes.Buffer{}
			if buffer.Cap() < int(fileEntry.Size) {
				buffer.Grow(int(fileEntry.Size) - buffer.Cap() + 1)
			}
			written, err := io.Copy(buffer, reader)
			if err != nil {
				b.Log("Cannot read file", fileEntry.Name, err)
			} else if written == 0 {
				b.Log("Empty file", fileEntry.Name)
			}
		}
		reader.Close()
		e = e.Next()
		if e == nil {
			e = fileQueue.Front()
		}
	}
}

func BenchmarkReadFromBufIOToPoolOfBuffer(b *testing.B) {
	b.ReportAllocs()
	if fileQueue.Len() == 0 {
		b.Skip("No source file")
	}
	localBufferPool := sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	e := fileQueue.Front()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileEntry := e.Value.(FileEntry)
		reader, err := getFileReader(fileEntry.Name)
		if err != nil {
			b.Log("Cannot open file", fileEntry.Name, err)
		} else {
			buffer := localBufferPool.Get().(*bytes.Buffer)
			buffer.Reset()
			if buffer.Cap() < int(fileEntry.Size) {
				buffer.Grow(int(fileEntry.Size) - buffer.Cap() + 1)
			}
			written, err := io.Copy(buffer, reader)
			if err != nil {
				b.Log("Cannot read file", fileEntry.Name, err)
			} else if written == 0 {
				b.Log("Empty file", fileEntry.Name)
			}
			localBufferPool.Put(buffer)
		}
		reader.Close()
		e = e.Next()
		if e == nil {
			e = fileQueue.Front()
		}
	}
}

func BenchmarkReadFromBufIOFromPoolToBuffer(b *testing.B) {
	b.ReportAllocs()
	if fileQueue.Len() == 0 {
		b.Skip("No source file")
	}
	e := fileQueue.Front()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileEntry := e.Value.(FileEntry)
		reader, err := getFileReaderFromPool(fileEntry.Name)
		if err != nil {
			b.Log("Cannot open file", fileEntry.Name, err)
		} else {
			buffer := &bytes.Buffer{}
			if buffer.Cap() < int(fileEntry.Size) {
				buffer.Grow(int(fileEntry.Size) - buffer.Cap() + 1)
			}
			written, err := io.Copy(buffer, reader)
			if err != nil {
				b.Log("Cannot read file", fileEntry.Name, err)
			} else if written == 0 {
				b.Log("Empty file", fileEntry.Name)
			}
		}
		reader.Close()
		bufReaderPool.Put(reader)
		e = e.Next()
		if e == nil {
			e = fileQueue.Front()
		}
	}
}

func BenchmarkReadFromBufIOFromPoolToPoolOfBuffer(b *testing.B) {
	b.ReportAllocs()
	if fileQueue.Len() == 0 {
		b.Skip("No source file")
	}
	localBufferPool := sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	e := fileQueue.Front()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileEntry := e.Value.(FileEntry)
		reader, err := getFileReaderFromPool(fileEntry.Name)
		if err != nil {
			b.Log("Cannot open file", fileEntry.Name, err)
		} else {
			buffer := localBufferPool.Get().(*bytes.Buffer)
			buffer.Reset()
			if buffer.Cap() < int(fileEntry.Size) {
				buffer.Grow(int(fileEntry.Size) - buffer.Cap() + 1)
			}
			written, err := io.Copy(buffer, reader)
			if err != nil {
				b.Log("Cannot read file", fileEntry.Name, err)
			} else if written == 0 {
				b.Log("Empty file", fileEntry.Name)
			}
			localBufferPool.Put(buffer)
		}
		reader.Close()
		bufReaderPool.Put(reader)
		e = e.Next()
		if e == nil {
			e = fileQueue.Front()
		}
	}
}

func BenchmarkReadFromBufIOToByteSlice(b *testing.B) {
	b.ReportAllocs()
	if fileQueue.Len() == 0 {
		b.Skip("No source file")
	}
	e := fileQueue.Front()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileEntry := e.Value.(FileEntry)
		reader, err := getFileReader(fileEntry.Name)
		if err != nil {
			b.Log("Cannot open file", fileEntry.Name, err)
		} else {
			buffer := make([]byte, 0, fileEntry.Size)
			read, err := reader.Read(buffer)
			if err != nil {
				b.Log("Cannot read file", fileEntry.Name, err)
			} else if read == 0 {
				b.Log("Empty file", fileEntry.Name)
			}
		}
		reader.Close()
		e = e.Next()
		if e == nil {
			e = fileQueue.Front()
		}
	}
}

func BenchmarkReadFromBufIOFromPoolToByteSlice(b *testing.B) {
	b.ReportAllocs()
	if fileQueue.Len() == 0 {
		b.Skip("No source file")
	}
	e := fileQueue.Front()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileEntry := e.Value.(FileEntry)
		reader, err := getFileReaderFromPool(fileEntry.Name)
		if err != nil {
			b.Log("Cannot open file", fileEntry.Name, err)
		} else {
			buffer := make([]byte, 0, fileEntry.Size)
			read, err := reader.Read(buffer)
			if err != nil {
				b.Log("Cannot read file", fileEntry.Name, err)
			} else if read == 0 {
				b.Log("Empty file", fileEntry.Name)
			}
		}
		reader.Close()
		bufReaderPool.Put(reader)
		e = e.Next()
		if e == nil {
			e = fileQueue.Front()
		}
	}
}
