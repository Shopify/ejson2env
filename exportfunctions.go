package ejson2env

import (
	"fmt"
	"io"
	"strings"

	"github.com/taskcluster/shell"
)

// ExportEnv writes the passed environment values to the passed
// io.Writer.
func ExportEnv(w io.Writer, values map[string]string) {
	for key, value := range values {
		fmt.Fprintf(w, "export %s=%s\n", key, shell.Escape(value))
	}
}

// ExportQuiet writes the passed environment values to the passed
// io.Writer in %s=%s format.
func ExportQuiet(w io.Writer, values map[string]string) {
	for key, value := range values {
		fmt.Fprintf(w, "%s=%s\n", key, shell.Escape(value))
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
