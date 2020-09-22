package main

import (
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
	"github.com/creativeprojects/clog"
)

type exclusion struct {
	globs     []string
	filenames []string
}

func newExclusion(patterns ...string) *exclusion {
	globs := []string{}
	filenames := []string{}
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}
		if strings.ContainsAny(pattern, `*/\`) {
			// This is a glob matching
			globs = append(globs, pattern)
		} else {
			// This is a filename matching only
			filenames = append(filenames, pattern)
		}
	}
	return &exclusion{
		globs:     globs,
		filenames: filenames,
	}
}

func (e *exclusion) match(fullname string) bool {
	return e.matchFilename(fullname) || e.matchPath(fullname)
}

func (e *exclusion) matchFilename(fullname string) bool {
	filename := filepath.Base(fullname)
	for _, pattern := range e.filenames {
		if filename == pattern {
			return true
		}
	}
	return false
}

func (e *exclusion) matchPath(fullname string) bool {
	for _, pattern := range e.globs {
		match, err := doublestar.PathMatch(pattern, fullname)
		if err != nil {
			clog.Error("invalid pattern")
		}
		if match {
			return true
		}
	}
	return false
}
