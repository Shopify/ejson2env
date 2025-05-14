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
// io.Writer.
func ExportEnv(w io.Writer, values map[string]string) {
	export(w, "export ", values)
}

// ExportQuiet writes the passed environment values to the passed
// io.Writer in %s=%s format.
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
			fmt.Fprintf(os.Stderr, "ejson2env blocked invalid key")
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

func validKey(k string) bool {
	for _, r := range k {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return false
		}
	}
	return true
}

func filteredValue(v string) string {
	printable := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, v)

	if printable != v {
		fmt.Fprintf(os.Stderr, "ejson2env trimmed control characters from value")
	}

	return shellescape.Quote(printable)
}
