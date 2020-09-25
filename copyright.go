package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

const (
	magicYear  = "###Year#From#Template#Here###"
	yearRegexp = `([\d]{4})`
)

// CopyrightData contains copyright template data.
// maybe we can add fields coming from the YAML configuration file?
type CopyrightData struct {
	Year int
}

type CopyrightTemplate struct {
	tmpl *template.Template
}

func ParseCopyrightTemplateFromFile(filename string) (*CopyrightTemplate, error) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		return nil, err
	}
	return &CopyrightTemplate{
		tmpl: tmpl,
	}, nil
}

func ParseCopyrightTemplateFromString(raw string) (*CopyrightTemplate, error) {
	tmpl := template.New("copyright")
	_, err := tmpl.Parse(raw)
	if err != nil {
		return nil, err
	}
	return &CopyrightTemplate{
		tmpl: tmpl,
	}, nil
}

// GetCopyrightNotice returns the copyright header with the template variables replaced by their values.
// please note any UTF8 BOM at the start of the template is stripped from the output.
func (t *CopyrightTemplate) GetCopyrightNotice(data interface{}) ([]byte, error) {
	// also use default buffer size to avoid unnecessary memory allocations
	buffer := bytes.NewBuffer(make([]byte, 0, defaultBufferSize))
	err := t.tmpl.Execute(buffer, &data)
	if err != nil {
		return nil, err
	}
	content := buffer.Bytes()
	if hasUTF8BOM(content) {
		return content[3:], nil
	}
	return content, nil
}

// getTextWithMagicValues will return the template filled in with magic values
//
// The magic values does not contain any special regexp character
func (t *CopyrightTemplate) getTextWithMagicValues() (string, error) {

	fakeData := map[string]string{
		"Year": magicYear,
	}
	// also use default buffer size to avoid unnecessary memory allocations
	buffer := bytes.NewBuffer(make([]byte, 0, defaultBufferSize))
	err := t.tmpl.Execute(buffer, &fakeData)
	if err != nil {
		return "", err
	}
	content := buffer.Bytes()
	if hasUTF8BOM(content) {
		return string(content[3:]), nil
	}
	return string(content), nil
}

// GetRegexp returns a searchable version of the template
func (t *CopyrightTemplate) GetRegexp() (*regexp.Regexp, error) {
	pattern, err := t.getTextWithMagicValues()
	if err != nil {
		return nil, err
	}
	return convertTextToRegexp(pattern)
}

func convertTextToRegexp(text string) (*regexp.Regexp, error) {
	escapeChars := []string{`\`, `^`, `$`, `.`, `|`, `?`, `*`, `+`, `(`, `)`, `[`, `]`, `{`, `}`}
	for _, escapeChar := range escapeChars {
		text = strings.ReplaceAll(text, escapeChar, `\`+escapeChar)
	}
	// put back any year into the template
	text = strings.ReplaceAll(text, magicYear, yearRegexp)
	// replace beginning of line by something more permissive
	text = strings.ReplaceAll(text, "\n \\*", "\n[ \t]*\\*")
	// replace end of line by something a bit more permissive
	text = strings.ReplaceAll(text, "\n", `[\s]+`)
	text = strings.ReplaceAll(text, "\r\n", `[\s]+`)
	// quick hack for the case a file only has a header with no return at the end
	if text[len(text)-1] == '+' {
		text = text[:len(text)-1] + "*"
	}
	fmt.Println(text)
	return regexp.Compile(text)
}
