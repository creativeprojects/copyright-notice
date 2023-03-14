package main

var (
	// UTF8BOM represents the 3 bytes of the BOM added by Microsoft IDEs
	UTF8BOM = []byte{0xef, 0xbb, 0xbf}
)

func hasUTF8BOM(content []byte) bool {
	if len(content) < 3 {
		return false
	}
	return content[0] == UTF8BOM[0] &&
		content[1] == UTF8BOM[1] &&
		content[2] == UTF8BOM[2]
}
