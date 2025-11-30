package helpproc

import (
	"regexp"
	"strings"
)

// NotesRegex matches a line that is exactly "Notes"
var NotesRegex = regexp.MustCompile(`^Notes$`)

// ProcessHelpOutput strips everything after a line matching ^Notes$ (inclusive)
func ProcessHelpOutput(output string) string {
	lines := strings.Split(output, "\n")
	var result []string
	for _, line := range lines {
		if NotesRegex.MatchString(line) {
			break
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// FormatHelpPreamble formats the preamble: TrimRight(preamble) + "\n\n" if non-empty
func FormatHelpPreamble(preamble string) string {
	trimmed := strings.TrimRight(preamble, " \t\n\r")
	if trimmed == "" {
		return ""
	}
	return trimmed + "\n\n"
}
