package ejson2env

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/taskcluster/shell"
)

var exportCommandPattern = regexp.MustCompile(`^export [a-zA-Z_][a-zA-Z0-9_-]*=([A-Za-z0-9/:=-]+|'.*')$`)
var quietCommandPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*=([A-Za-z0-9/:=-]+|'.*')$`)

// ValidateExportCommand checks if a string matches the format: export some_username='anycharacters'
func ValidateExportCommand(cmd string) bool {
	return exportCommandPattern.MatchString(cmd)
}

// ValidateQuietCommand checks if a string matches the format: some_username='anycharacters'
func ValidateQuietCommand(cmd string) bool {
	return quietCommandPattern.MatchString(cmd)
}

// ExportEnv writes the passed environment values to the passed
// io.Writer.
func ExportEnv(w io.Writer, values map[string]string) {
	for key, value := range values {
		cmd := fmt.Sprintf("export %s=%s", key, escape(value))
		if ValidateExportCommand(cmd) {
			fmt.Fprintln(w, cmd)
		}
	}
}

// ExportQuiet writes the passed environment values to the passed
// io.Writer in %s=%s format.
func ExportQuiet(w io.Writer, values map[string]string) {
	for key, value := range values {
		cmd := fmt.Sprintf("%s=%s", key, escape(value))
		if ValidateQuietCommand(cmd) {
			fmt.Fprintln(w, cmd)
		}
	}
}

func TrimLeadingUnderscoreExportWrapper(exportfunc ExportFunction) ExportFunction {
	return func(w io.Writer, values map[string]string) {
		newValues := make(map[string]string, len(values))

		for key, value := range values {
			newValues[strings.TrimLeft(key, "_")] = value
		}

		exportfunc(w, newValues)
	}
}

func escape(v string) string {
	printable := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, v)
	return shell.Escape(printable)
}
