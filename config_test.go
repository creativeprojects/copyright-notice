package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	source := `---
profiles:
  self:
    source: .
    extensions: go
    utf8-bom: forget
    year: update
    exclude-from: exclude
    copyright: copyright.txt
    detect-own: '^([ \t\r\n]*/\*[ \t\r\n*]*Copyright[ \t]+\(C\)[ \t]+)([\d]{4})([ \t]+CreativeProjects\.[ \t]*[\r\n]+)'
    detect-others: '^[ \t\r\n]*/\*[ \t\r\n*]*Copyright[ \t]+'

  gopath:
    source: ~/go
    extensions:
      - go
      - js
    utf8-bom: keep
    excludes:
      - node_modules
      - "**/.*"
    copyright: short-copyright.txt
`
	config, err := LoadConfig(bytes.NewBufferString(source))
	require.NoError(t, err)
	assert.NotEmpty(t, config)
}

func TestCleanupFileExtensions(t *testing.T) {
	config := Config{
		Profiles: map[string]ConfigProfile{
			"first": {
				Extensions: &StringSlice{"js", ".ts"},
			},
			"second": {
				Extensions: &StringSlice{".cs", "yaml"},
			},
		},
	}
	cleanupConfig(&config)
	assert.Len(t, config.Profiles, 2)
	assert.ElementsMatch(t, []string(*config.Profiles["first"].Extensions), []string{".js", ".ts"})
	assert.ElementsMatch(t, []string(*config.Profiles["second"].Extensions), []string{".cs", ".yaml"})
}
