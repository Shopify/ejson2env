package ejson2env

import (
	"bytes"
	"strings"
	"testing"
)

func TestValidKeyAcceptsOnlyShellNames(t *testing.T) {
	for key, want := range map[string]bool{
		"A":      true,
		"_a":     true,
		"DB_1":   true,
		"":       false,
		"1A":     false,
		"A-B":    false,
		"A B":    false,
		"=":      false,
		"é":      false,
		"A\nB":   false,
		"a.b":    false,
		"_":      true, // valid on its own; --trim-underscore turns it into ""
		"__DB__": true,
	} {
		if got := validKey(key); got != want {
			t.Errorf("validKey(%q) = %v, want %v", key, got, want)
		}
	}
}

// --trim-underscore rewrites keys after ExtractEnv validated them, so export must
// re-validate. A key of "_" becomes "" and would print a line starting with "=",
// which zsh executes as a command.
func TestTrimUnderscoreDropsEmptyAndInvalidKeys(t *testing.T) {
	values := map[string]string{
		"_":    "/tmp/not-a-command",
		"_1A":  "leading digit after trim",
		"_OK":  "kept",
		"SAFE": "kept",
	}
	for name, export := range map[string]ExportFunction{
		"ExportEnv": ExportEnv, "ExportQuiet": ExportQuiet,
	} {
		var out bytes.Buffer
		TrimLeadingUnderscoreExportWrapper(export)(&out, values)
		lines := strings.Split(strings.TrimSpace(out.String()), "\n")
		for _, line := range lines {
			assignment := strings.TrimPrefix(line, "export ")
			if !strings.HasPrefix(assignment, "OK=") && !strings.HasPrefix(assignment, "SAFE=") {
				t.Errorf("%s emitted an unexpected line %q", name, line)
			}
		}
		if len(lines) != 2 {
			t.Errorf("%s emitted %d lines, want 2: %q", name, len(lines), out.String())
		}
	}
}

func TestFilteredValueKeepsTabsAndNewlines(t *testing.T) {
	in := "a\tb\nc\x00d\x1be"
	if got, want := filteredValue(in), "'a\tb\ncde'"; got != want {
		t.Errorf("filteredValue(%q) = %q, want %q", in, got, want)
	}
}
