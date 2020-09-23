package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	expectedCopyright = "/* Copyright 2020 TestCorp */\n"
)

func TestCopyrightWithoutBOM(t *testing.T) {
	data := copyrightData{Year: 2020}
	content, err := getCopyrightNoticeFromTemplate("test_files/copyright_without_BOM.js", &data)
	require.NoError(t, err)
	assert.Equal(t, expectedCopyright, string(content))
}

func TestCopyrightWithBOM(t *testing.T) {
	data := copyrightData{Year: 2020}
	content, err := getCopyrightNoticeFromTemplate("test_files/copyright_with_BOM.js", &data)
	require.NoError(t, err)
	assert.Equal(t, expectedCopyright, string(content))
}
