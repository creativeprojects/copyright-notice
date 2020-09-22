package main

type fileStatus uint8

const (
	fileStatusUnknown fileStatus = iota
	fileStatusNoCopyright
	fileStatusWithCopyright
	fileStatusCopyrightYearNeedsUpdated
	fileStatusIgnore
	fileStatusTooBig
	fileStatusCannotOpen
	fileStatusError // Keep this one last!
)

func (f fileStatus) String() string {
	switch f {
	case fileStatusNoCopyright:
		return "add copyright data"
	case fileStatusWithCopyright:
		return "found copyright data"
	case fileStatusCopyrightYearNeedsUpdated:
		return "copyright year needs updated"
	case fileStatusIgnore:
		return "ignore auto-generated file"
	case fileStatusTooBig:
		return "file is too big"
	case fileStatusCannotOpen:
		return "cannot open file"
	case fileStatusError:
		return "general read/write error"
	}
	return ""
}

func (f fileStatus) Symbol() string {
	switch f {
	case fileStatusNoCopyright:
		return "+"
	case fileStatusWithCopyright:
		return "."
	case fileStatusCopyrightYearNeedsUpdated:
		return "^"
	case fileStatusIgnore:
		return "-"
	case fileStatusTooBig:
		return "O"
	case fileStatusError:
		return "!"
	case fileStatusCannotOpen:
		return "X"
	case fileStatusUnknown:
		return "?"
	}
	return " "
}
