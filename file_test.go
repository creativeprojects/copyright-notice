package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileErrorInvalidDescriptor(t *testing.T) {
	file := NewFile(bufferSize)
	require.NotNil(t, file)

	// run the test twice
	for i := 0; i < 2; i++ {
		err := file.Read("", bufferSize)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), FileErrorInvalidName)
		}

		// reuse the same buffer for the next file in test
		file.Reset()
	}
}

func TestFileErrorTooBig(t *testing.T) {
	size := 10
	file := NewFile(size)
	require.NotNil(t, file)

	// run the test twice
	for i := 0; i < 2; i++ {
		err := file.Read("some file.txt", int64(size+1))
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), FileErrorTooBig)
		}
		assert.False(t, file.IsReady())

		// reuse the same buffer for the next file in test
		file.Reset()
	}
}

func TestFileErrorCannotOpen(t *testing.T) {
	file := NewFile(bufferSize)
	require.NotNil(t, file)

	// run the test twice
	for i := 0; i < 2; i++ {
		err := file.Read("some file.txt", 10)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "cannot open file")
		}
		assert.False(t, file.IsReady())

		// reuse the same buffer for the next file in test
		file.Reset()
	}
}

func TestFileReadTooLittle(t *testing.T) {
	var size int64 = 1024
	file := NewFile(bufferSize)
	require.NotNil(t, file)

	// run the test twice
	for i := 0; i < 2; i++ {
		err := file.Read(fmt.Sprintf("test_files/random%d.txt", size), size+1)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), fmt.Sprintf("file size = %d bytes but read %d bytes instead", size+1, size))
		}
		assert.False(t, file.IsReady())

		// reuse the same buffer for the next file in test
		file.Reset()
	}
}

func TestFileReadWithoutBOM(t *testing.T) {
	name := "test_files/without_BOM.txt"
	info, err := os.Stat(name)
	require.NoError(t, err)

	file := NewFile(bufferSize)
	require.NotNil(t, file)

	// run the test twice
	for i := 0; i < 2; i++ {
		err := file.Read(name, info.Size())
		assert.NoError(t, err)
		assert.True(t, file.IsReady())
		assert.False(t, file.HasUTF8BOM())
		// first byte should be an 'f'
		assert.Equal(t, byte('f'), file.Bytes()[0])

		// reuse the same buffer for the next file in test
		file.Reset()
	}
}

func TestFileReadWithBOM(t *testing.T) {
	name := "test_files/with_BOM.txt"
	info, err := os.Stat(name)
	require.NoError(t, err)

	file := NewFile(bufferSize)
	require.NotNil(t, file)

	// run the test twice
	for i := 0; i < 2; i++ {
		err := file.Read(name, info.Size())
		assert.NoError(t, err)
		assert.True(t, file.IsReady())
		assert.True(t, file.HasUTF8BOM())
		// first byte should be an 'f'
		assert.Equal(t, byte('f'), file.Bytes()[0])

		// reuse the same buffer for the next file in test
		file.Reset()
	}
}
