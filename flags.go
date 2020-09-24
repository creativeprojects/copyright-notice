package main

import (
	flag "github.com/spf13/pflag"
)

type Flags struct {
	verbose        bool
	dryRun         bool
	configFile     string
	outputFilename string
	help           bool
}

var (
	flags = &Flags{}
)

func init() {
	flag.StringVarP(&flags.configFile, "config", "c", "copyright-notice.yaml", "Configuration file")
	flag.BoolVar(&flags.dryRun, "dry-run", false, "Show all the files that would be processed, but don't save anything")
	flag.BoolVarP(&flags.verbose, "verbose", "v", false, "Display more information")
	flag.StringVarP(&flags.outputFilename, "output", "o", "", "Write the output into a file instead of the console")
	flag.BoolVarP(&flags.help, "help", "h", false, "Prints usage")
}
