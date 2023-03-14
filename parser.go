package main

import (
	"container/list"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/creativeprojects/clog"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

type Parser struct {
	extensions []string
	exclusions *exclusion
	fileQueue  *list.List
}

func NewParser(extensions []string, exclusions *exclusion) *Parser {
	return &Parser{
		extensions: extensions,
		exclusions: exclusions,
		fileQueue:  list.New(),
	}
}

func (p *Parser) Directories(directories []string) *list.List {
	if directories == nil || len(directories) == 0 {
		return p.fileQueue
	}
	total := int64(len(directories))
	progress := mpb.New(nil)
	spinner := progress.AddSpinner(total, mpb.SpinnerOnLeft,
		mpb.PrependDecorators(decor.CountersNoUnit("directories and files analyzed: %d / %d", decor.WC{})),
		mpb.BarRemoveOnComplete(),
	)
	for _, source := range directories {
		if source == "" {
			continue
		}
		p.directory(source,
			func(more int) {
				total += int64(more)
				spinner.SetTotal(total, false)
			},
			func() {
				spinner.Increment()
			})
	}
	spinner.SetTotal(total, true)
	progress.Wait()
	return p.fileQueue
}

func (p *Parser) directory(directory string, addTotal func(int), addFile func()) {
	directory = filepath.Clean(directory)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		clog.Errorf("cannot parse directory: %s", err)
	}
	if files == nil || len(files) == 0 {
		return
	}
	addTotal(len(files))
	for _, file := range files {
		// Make sure we don't go into a infinite loop when running on unixes
		if file.Name() == "." || file.Name() == ".." {
			continue
		}

		fullName := filepath.Join(directory, file.Name())
		if p.exclusions.match(fullName) {
			clog.Debugf("path excluded: '%s'", fullName)
			continue
		}
		addFile()
		if file.IsDir() {
			p.directory(fullName, addTotal, addFile)
		} else if file.Size() > minFileSize && p.matchExtension(file.Name()) {
			p.fileQueue.PushBack(FileEntry{fullName, file.Size()})
			// update the max size of the files we're going to analyze,
			// we keep the oversized file for reporting, but we won't build
			// a buffer of that size
			if file.Size() > maxSize && file.Size() <= maxFileSize {
				maxSize = file.Size()
			}
		}
	}
}

func (p *Parser) matchExtension(fileName string) bool {
	for _, ext := range p.extensions {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}
