package main

import (
	flag "github.com/spf13/pflag"
)

type Flags struct {
	verbose           bool
	dryRun            bool
	configFile        string
	sourceDirectory   []string
	copyrightFilename string
	excludeFilename   string
	outputFilename    string
	exclude           []string
	extensions        []string
	help              bool
}

var (
	flags = &Flags{}
)

func init() {
	flag.BoolVarP(&flags.verbose, "verbose", "v", false, "Display more information")
	flag.BoolVar(&flags.dryRun, "dry-run", false, "Show all the files that would be updated, but don't save anything")
	flag.StringVarP(&flags.configFile, "config", "c", "", "Configuration file")
	flag.StringSliceVarP(&flags.sourceDirectory, "source", "s", []string{"."}, "Directory containing your source files: can be specified more than once")
	flag.StringVarP(&flags.copyrightFilename, "copyright", "t", "copyright.txt", "File containing the copyright notice template")
	flag.StringVar(&flags.excludeFilename, "exclude-from", "exclude", "File containing a list of exclusion patterns")
	flag.StringVarP(&flags.outputFilename, "output", "o", "", "Write the output into a file instead of the console")
	flag.StringSliceVar(&flags.exclude, "exclude", nil, "Exclude pattern: can be specified more than once")
	flag.StringSliceVar(&flags.extensions, "ext", nil, "File extensions to include: can be specified more than once")
	flag.BoolVarP(&flags.help, "help", "h", false, "Prints usage")
}

func cleanupConfiguration() {
	// We'll prefix the files extension by a dot if it's not there yet
	for index, extension := range flags.extensions {
		if extension[0] != '.' {
			extension = "." + extension
		}
		flags.extensions[index] = extension
	}
}
