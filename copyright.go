package main

import (
	"bytes"
	"text/template"
)

type copyrightData struct {
	Year int
}

func getFullCopyrightNoticeFromTemplate(data interface{}) ([]byte, error) {
	return loadTemplate(config.copyrightFilename, data)
}

func loadTemplate(filename string, data interface{}) ([]byte, error) {
	copyrightTemplate, err := template.ParseFiles(filename)
	if err != nil {
		return nil, err
	}
	buffer := &bytes.Buffer{}
	err = copyrightTemplate.Execute(buffer, &data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
