package ejson2env

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"

	"al.essio.dev/pkg/shellescape"
)

// ExportEnv writes the passed environment values to the passed
// io.Writer as POSIX shell code. When using eval, pass the output as a single,
// double-quoted argument to prevent word splitting and filename expansion.
func ExportEnv(w io.Writer, values map[string]string) {
	export(w, "export ", values)
}

// ExportQuiet writes the passed environment values to the passed
// io.Writer in %s=%s format, with the same shell quoting requirements as ExportEnv.
func ExportQuiet(w io.Writer, values map[string]string) {
	export(w, "", values)
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

func export(w io.Writer, prefix string, values map[string]string) {
	keys := make([]string, 0, len(values))
	for k := range values {
		if !validKey(k) {
			fmt.Fprintln(os.Stderr, "ejson2env blocked invalid key")
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		value := filteredValue(values[k])
		fmt.Fprintf(w, "%s%s=%s\n", prefix, k, value)
	}
}

// validKey accepts only POSIX shell names. Keys are validated when the file is
// read, but options such as --trim-underscore rewrite them afterwards, and an
// empty or otherwise invalid name here would change what the shell executes.
func validKey(k string) bool {
	return validIdentifierPattern.MatchString(k)
}

func filteredValue(v string) string {
	printable := strings.Map(func(r rune) rune {
		// Tabs and newlines are safe inside the single quotes the exporter emits.
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, v)

	if printable != v {
		fmt.Fprintln(os.Stderr, "ejson2env trimmed control characters from value")
	}

	quoted := shellescape.Quote(printable)
	// shellescape leaves '=' unquoted, but zsh expands a leading '=' in an
	// assignment value to the corresponding command path.
	if strings.HasPrefix(quoted, "=") {
		return "'" + quoted + "'"
	}
	return quoted
}
