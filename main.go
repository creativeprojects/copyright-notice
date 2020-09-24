package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/creativeprojects/clog"
	flag "github.com/spf13/pflag"
)

const (
	minFileSize = 3
	maxFileSize = 2 * 1024 * 1024
)

type resultData struct {
	fileName string
	err      error
}

var (
	results           []*list.List
	detectOwnHeader   *regexp.Regexp
	detectOtherHeader *regexp.Regexp
	autoGenerated     *regexp.Regexp
	randomGenerator   *rand.Rand
	maxSize           int64
)

func init() {
	results = make([]*list.List, fileStatusError+1)
	for i := 0; i <= int(fileStatusError); i++ {
		results[i] = &list.List{}
	}
	randomGenerator = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	/* regexp explanation:
	combination of end of line, spaces and tabs (only)
	start of comment: "/*"
	combination of end of line, spaces, tabs and * (only)
	word "Copyright"
	spaces or tabs (only)
	string "(C)"
	spaces or tabs (only)
	4 digits
	spaces or tabs (only)
	string "CreativeProjects."
	spaces or tabs (only)
	end of line (only)
	... at that point we don't need to check any further
	*/
	detectOwnHeader = regexp.MustCompile(`^([ \t\r\n]*/\*[ \t\r\n*]*Copyright[ \t]+\(C\)[ \t]+)([\d]{4})([ \t]+CreativeProjects\.[ \t]*[\r\n]+)`)
	detectOtherHeader = regexp.MustCompile(`^[ \t\r\n]*/\*[ \t\r\n*]*Copyright[ \t]+`)
	autoGenerated = regexp.MustCompile(`\<auto-generated\>`)
}

func main() {

	flag.Parse()
	if flags.help {
		fmt.Print("\nUsage of copyright-notice:\n\n")
		flag.PrintDefaults()
		return
	}
	close := setupLogger(flags)
	defer close()

	// load configuration
	config, err := LoadFileConfig(flags.configFile)
	if err != nil {
		clog.Errorf("cannot open configuration file: %s", err)
	}

	for name, profile := range config.Profiles {
		// log prefix should be displayed only if we have more than one profile
		if len(config.Profiles) > 1 {
			clog.SetPrefix(name + ":  ")
		}
		if profile.Source == nil || len(*profile.Source) == 0 {
			clog.Warning("no source folder defined, skipping profile")
			continue
		}
		if profile.Extensions == nil || len(*profile.Extensions) == 0 {
			clog.Warning("no file extension defined, skipping profile")
			continue
		}
		if profile.Copyright == "" {
			clog.Warning("no copyright file defined, skipping profile")
			continue
		}
		clog.Infof("searching for source files %s in folder %s", *profile.Extensions, *profile.Source)

		var excludeList []string
		// Load exclusion list from file
		if profile.ExcludeFrom != "" {
			excludeList, err = readLines(profile.ExcludeFrom)
			if err != nil {
				clog.Warningf("error while reading exclusion file: %s, skipping profile", err)
				continue
			}
		}
		if profile.Excludes != nil && len(*profile.Excludes) > 0 {
			excludeList = append(excludeList, *profile.Excludes...)
		}
		// Generate the exclusions
		exclusions := newExclusion(excludeList...)

		// Parse the source directory for files
		parser := NewParser(*profile.Extensions, exclusions)
		fileQueue := parser.Directories(*profile.Source)
		if fileQueue.Len() == 0 {
			clog.Warning("no matching file found")
			continue
		}

		// Load the copyright notice template
		copyrightNotice, err := getCopyrightNoticeFromTemplate(profile.Copyright, &copyrightData{Year: time.Now().Year()})
		if err != nil {
			clog.Errorf("cannot load copyright template: %v", err)
			return
		}

		// Merge all files with the copyright notice
		clog.Infof("analyzing %d source files", fileQueue.Len())
		checkForCopyrightNotices(fileQueue, copyrightNotice)
		// fmt.Println("")

		// Display results in debug mode
		if flags.verbose {
			displayDetailedResults()
		} else {
			displaySummaryResults()
		}
		clog.SetPrefix("")
	}

}

func progress(fileName string, status fileStatus, err error) {
	// Keep results for later use
	results[status].PushBack(&resultData{fileName, err})
	// Display live progress
	// fmt.Print(status.Symbol())
}

func displayDetailedResults() {
	for _, status := range []fileStatus{
		fileStatusNoCopyright,
		fileStatusWithCopyright,
		fileStatusCopyrightYearNeedsUpdated,
		fileStatusIgnore,
		fileStatusTooBig,
		fileStatusCannotOpen,
		fileStatusError,
	} {
		displayResultList(results[status], status.String())
	}
}

func displayResultList(list *list.List, statusMessage string) {
	if list == nil || list.Len() == 0 {
		return
	}
	for e := list.Front(); e != nil; e = e.Next() {
		status := e.Value.(*resultData)
		details := fmt.Sprintf(", file: '%s'", status.fileName)
		if status.err != nil {
			details += fmt.Sprintf(", error: %s", status.err)
		}
		clog.Debug(statusMessage + details)
	}
}

func displaySummaryResults() {
	for _, status := range []fileStatus{
		fileStatusNoCopyright,
		fileStatusWithCopyright,
		fileStatusCopyrightYearNeedsUpdated,
		fileStatusIgnore,
		fileStatusTooBig,
		fileStatusCannotOpen,
		fileStatusError,
	} {
		displaySummary(results[status], status.String())
	}
}

func displaySummary(list *list.List, message string) {
	if list == nil || list.Len() == 0 {
		return
	}
	clog.Infof(message+": %d %s", list.Len(), simplePlural("file", list.Len()))
}

func simplePlural(word string, count int) string {
	if count > 1 {
		return word + "s"
	}
	return word
}
