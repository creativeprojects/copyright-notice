package main

import (
	"io"

	"gopkg.in/yaml.v2"
)

// StringSlice is used to allow a YAML field to be represented by a single string or an array of strings
//
//   ---
//   works: "single"
//   also-works:
//     - "one"
//     - "two"
type StringSlice []string

// UnmarshalYAML a string or a list of strings into a slice
func (a *StringSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

type Config struct {
	Profiles map[string]ConfigProfiles `yaml:"profiles"`
}

type ConfigProfiles struct {
	Source       *StringSlice `yaml:"source"`
	Extensions   *StringSlice `yaml:"extensions"`
	BOM          string       `yaml:"utf8-bom"`
	Year         string       `yaml:"year"`
	ExcludeFrom  string       `yaml:"exclude-from"`
	Excludes     *StringSlice `yaml:"excludes"`
	Copyright    string       `yaml:"copyright"`
	DetectOwn    string       `yaml:"detect-own"`
	DetectOthers string       `yaml:"detect-others"`
}

func LoadConfig(reader io.Reader) (Config, error) {
	decoder := yaml.NewDecoder(reader)
	config := Config{}
	err := decoder.Decode(&config)
	return config, err
}
