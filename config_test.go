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
