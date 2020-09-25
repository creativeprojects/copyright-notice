package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	expectedCopyright = "/* Copyright 2020 TestCorp */\n"
)

func TestCopyrightWithoutBOM(t *testing.T) {
	data := CopyrightData{Year: 2020}
	tmpl, err := ParseCopyrightTemplateFromFile("test_files/copyright_without_BOM.js")
	require.NoError(t, err)
	content, err := tmpl.GetCopyrightNotice(&data)
	require.NoError(t, err)
	assert.Equal(t, expectedCopyright, string(content))
}

func TestCopyrightWithBOM(t *testing.T) {
	data := CopyrightData{Year: 2020}
	tmpl, err := ParseCopyrightTemplateFromFile("test_files/copyright_with_BOM.js")
	require.NoError(t, err)
	content, err := tmpl.GetCopyrightNotice(&data)
	require.NoError(t, err)
	assert.Equal(t, expectedCopyright, string(content))
}

func TestConversionToRegexp(t *testing.T) {
	testData := []string{
		"",
		"Copyright",
		"\ncopyright \n",
		"\\Copyright",
		"Copyright^",
		"$Copyright",
		".Copyright.",
		"Copyright|",
		"Copyright?",
		"*Copyright",
		"+Copyright",
		")Copyright(",
		"[copyright]",
		"}copyright{",
	}
	for index, testItem := range testData {
		text := testItem
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			pattern, err := convertTextToRegexp(text)
			require.NoError(t, err)
			result := pattern.FindString(text)
			assert.Equal(t, text, result)
		})
	}
}

func TestTemporaryTemplateForSearching(t *testing.T) {
	raw := "Copyright (c) {{ .Year }} *Some+Corp?"
	tmpl, err := ParseCopyrightTemplateFromString(raw)
	require.NoError(t, err)
	text, err := tmpl.getTextWithMagicValues()
	require.NoError(t, err)
	assert.Equal(t, strings.Replace(raw, "{{ .Year }}", magicYear, 1), text)
}
