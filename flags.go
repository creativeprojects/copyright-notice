package main

import (
	"bufio"
	"os"

	flag "github.com/spf13/pflag"
)

type flags struct {
	verbose            bool
	dryRun             bool
	confirmNoExclusion bool
	sourceDirectory    []string
	copyrightFilename  string
	excludeFilename    string
	outputFilename     string
	exclude            []string
	extensions         []string
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
// this is used to read an exclusion file
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

var (
	config = &flags{}
)

func init() {
	flag.BoolVarP(&config.verbose, "verbose", "v", false, "Display more information")
	flag.BoolVar(&config.dryRun, "dry-run", false, "Show all the files that would be updated, but don't save anything")
	flag.BoolVarP(&config.confirmNoExclusion, "empty-exclude", "e", false, "Confirm this is ok to have an exclusion list empty. Please note it does not empty the list if there's one")
	flag.StringSliceVarP(&config.sourceDirectory, "source", "s", []string{"."}, "Directory containing your source files: can be specified more than once")
	flag.StringVarP(&config.copyrightFilename, "copyright", "c", "copyright.txt", "File containing the copyright notice template")
	flag.StringVar(&config.excludeFilename, "exclude-from", "exclude", "File containing a list of exclusion patterns")
	flag.StringVarP(&config.outputFilename, "output", "o", "", "Write the output into a file instead of the console")
	flag.StringSliceVar(&config.exclude, "exclude", nil, "Exclude pattern: can be specified more than once")
	flag.StringSliceVar(&config.extensions, "ext", nil, "File extensions to include: can be specified more than once")
}

func cleanupConfiguration() {
	// We'll prefix the files extension by a dot if it's not there yet
	for index, extension := range config.extensions {
		if extension[0] != '.' {
			extension = "." + extension
		}
		config.extensions[index] = extension
	}
}
