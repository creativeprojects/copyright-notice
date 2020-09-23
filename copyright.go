package main

import (
	"bytes"
	"text/template"
)

// copyrightData contains copyright template data.
// maybe we can add fields coming from the YAML configuration file?
type copyrightData struct {
	Year int
}

// getCopyrightNoticeFromTemplate returns the copyright header with the template variables replaced by their values.
// please note any UTF8 BOM at the start of the template is stripped from the output.
func getCopyrightNoticeFromTemplate(filename string, data interface{}) ([]byte, error) {
	copyrightTemplate, err := template.ParseFiles(filename)
	if err != nil {
		return nil, err
	}
	// also use default buffer size to avoid unnecessary memory allocations
	buffer := bytes.NewBuffer(make([]byte, 0, defaultBufferSize))
	err = copyrightTemplate.Execute(buffer, &data)
	if err != nil {
		return nil, err
	}
	content := buffer.Bytes()
	if hasUTF8BOM(content) {
		return content[3:], nil
	}
	return content, nil
}
