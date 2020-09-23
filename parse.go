package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/creativeprojects/clog"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

func parseDirectories(directories []string, exclusions *exclusion) {
	if directories == nil || len(directories) == 0 {
		return
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
		parseDirectory(source, exclusions,
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
}

func parseDirectory(directory string, exclusions *exclusion, addTotal func(int), addFile func()) {
	directory = filepath.Clean(directory)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		clog.Errorf("cannot parse directory: %w", err)
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
		if exclusions.match(fullName) {
			clog.Debugf("path excluded: '%s'", fullName)
			continue
		}
		addFile()
		if file.IsDir() {
			parseDirectory(fullName, exclusions, addTotal, addFile)
		} else if file.Size() > minFileSize && matchExtension(file.Name()) {
			fileQueue.PushBack(FileEntry{fullName, file.Size()})
			if file.Size() > maxSize {
				maxSize = file.Size()
			}
		}
	}
}

func matchExtension(fileName string) bool {
	for _, ext := range flags.extensions {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}
