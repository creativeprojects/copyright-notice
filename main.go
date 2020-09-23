package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
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
	fileQueue         *list.List
	results           []*list.List
	detectOwnHeader   *regexp.Regexp
	detectOtherHeader *regexp.Regexp
	autoGenerated     *regexp.Regexp
	randomGenerator   *rand.Rand
	maxSize           int64
)

func init() {
	fileQueue = list.New()
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
	var err error

	flag.Parse()
	if flags.help {
		fmt.Print("\nUsage of copyright-notice:\n\n")
		flag.PrintDefaults()
		return
	}
	setupLogger()
	clog.Infof("starting copyright notice in source folder: '%s'", flags.sourceDirectory)

	// We need at least one file extension
	if len(flags.extensions) == 0 {
		clog.Error("you haven't specified any file extension, for example '--ext js' for Javascript files")
		return
	}
	cleanupConfiguration()
	clog.Infof("searching for source files with extensions %v", flags.extensions)

	// Load exclusion list from file
	var excludeList []string
	if flags.excludeFilename != "" {
		excludeList, err = readLines(flags.excludeFilename)
		if err != nil {
			clog.Errorf("error while reading exclusion file: %w", err)
		}
	}
	// Generate the exclusions
	exclusions := newExclusion(append(excludeList, flags.exclude...)...)

	// Parse the source directory for files
	parseDirectories(flags.sourceDirectory, exclusions)
	if fileQueue.Len() == 0 {
		clog.Warning("found absolutely no file matching these extensions")
		return
	}

	// Load the copyright notice template
	clog.Infof("analyzing %d source files", fileQueue.Len())
	copyrightNotice, err := getCopyrightNoticeFromTemplate(flags.copyrightFilename, &copyrightData{Year: time.Now().Year()})
	if err != nil {
		clog.Errorf("cannot load copyright template: %w", err)
		return
	}

	// Merge all files with the copyright notice
	checkForCopyrightNotices(copyrightNotice)
	fmt.Println("")

	// Display results in debug mode
	if flags.verbose {
		displayDetailedResults()
	} else {
		displaySummaryResults()
	}
}

func addCopyrightNotice(fileName string, copyrightNotice, buffer []byte) error {
	var err error
	randomBytes := make([]byte, 10)
	randomGenerator.Read(randomBytes)
	tempFilename := filepath.Join(filepath.Dir(fileName), "$"+fmt.Sprintf("%x", randomBytes)+"$"+filepath.Base(fileName))

	err = createFile(tempFilename, copyrightNotice, buffer, false)
	if err != nil {
		return err
	}
	// Move the temp file into place
	err = os.Rename(tempFilename, fileName)
	if err != nil {
		// Try to delete the temp file
		os.Remove(tempFilename)
		return err
	}
	return nil
}

func createFile(fileName string, header, content []byte, withBOM bool) error {
	outputFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write the copyright notice first
	_, err = outputFile.Write(header)
	if err != nil {
		return err
	}

	// Then write the file content
	outputFile.Write(content)
	if err != nil {
		return err
	}
	return nil
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
		details := fmt.Sprintf(". file: '%s'", status.fileName)
		if status.err != nil {
			details += fmt.Sprintf(". error: '%s'", status.err)
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
