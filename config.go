package main

import (
	"io"
	"os"

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
	MaxFileSize       int64                    `yaml:"max-file-size"`
	DefaultBufferSize int                      `yaml:"default-buffer-size"`
	Profiles          map[string]ConfigProfile `yaml:"profiles"`
}

type ConfigProfile struct {
	Source               *StringSlice `yaml:"source"`     // Mandatory
	Extensions           *StringSlice `yaml:"extensions"` // Mandatory
	Copyright            string       `yaml:"copyright"`  // Mandatory
	BOM                  string       `yaml:"utf8-bom"`
	Year                 string       `yaml:"year"`
	Excludes             *StringSlice `yaml:"excludes"`
	ExcludeFrom          string       `yaml:"exclude-from"`
	ExcludeFromGitIgnore string       `yaml:"exclude-gitignore"`
	DetectOwn            string       `yaml:"detect-own"`
	DetectOthers         string       `yaml:"detect-others"`
	CommitChanges        string       `yaml:"commit-changes"`
	CommitMessage        string       `yaml:"commit-message"`
	CommitAuthor         string       `yaml:"commit-author"`
	Output               string       `yaml:"output"`
}

// NewConfig creates a new configuration with the default values
func NewConfig() Config {
	return Config{
		MaxFileSize:       maxFileSize,
		DefaultBufferSize: defaultBufferSize,
	}
}

// LoadFileConfig loads YAML configuration from a file
func LoadFileConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return NewConfig(), err
	}
	defer file.Close()
	return LoadConfig(file)
}

// LoadConfig loads YAML configuration from a reader
func LoadConfig(reader io.Reader) (Config, error) {
	decoder := yaml.NewDecoder(reader)
	config := NewConfig()
	err := decoder.Decode(&config)
	cleanupConfig(&config)
	return config, err
}

func cleanupConfig(config *Config) {
	for _, profile := range config.Profiles {
		// we'll prefix the files extension by a dot if it's not there yet
		if profile.Extensions != nil {
			for index, extension := range *profile.Extensions {
				if extension[0] != '.' {
					extension = "." + extension
				}
				(*profile.Extensions)[index] = extension
			}
		}
		// we expend the environment variables in paths
		if profile.Source != nil {
			for index, dir := range *profile.Source {
				(*profile.Source)[index] = os.ExpandEnv(dir)
			}
		}
	}
}
