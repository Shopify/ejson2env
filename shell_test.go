package ejson2env_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Shopify/ejson2env/v2"
)

func testShell(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("requires a POSIX shell and POSIX filenames")
	}
	shell, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("requires a POSIX shell")
	}
	return shell
}

func TestExportShellValues(t *testing.T) {
	shell := testShell(t)
	cases := []struct{ name, value string }{
		{"empty", ""},
		{"wildcards", "* ? [abc]"},
		{"whitespace", "  two   words  "},
		{"quotes", "single ' and double \" and backslash \\"},
		{"commands", "$( : > marker) `: > marker`; : > marker; #"},
		{"newlines", "first\nlast\n\n"},
		{"tabs", "a\tb"},
		{"leading equals", "=true"},
	}
	for name, export := range map[string]ejson2env.ExportFunction{
		"ExportEnv": ejson2env.ExportEnv, "ExportQuiet": ejson2env.ExportQuiet,
	} {
		for _, tc := range cases {
			t.Run(name+"/"+tc.name, func(t *testing.T) {
				var output bytes.Buffer
				export(&output, map[string]string{"VALUE": tc.value})
				cmd := exec.Command(shell, "-c", `eval "$(printf '%s' "$OUTPUT")"; printf '%s' "$VALUE"`)
				cmd.Dir = t.TempDir()
				cmd.Env = []string{"OUTPUT=" + output.String(), "LC_ALL=C"}
				got, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("shell failed: %v\n%s", err, got)
				}
				if string(got) != tc.value {
					t.Errorf("got %q, want %q", got, tc.value)
				}
				if _, err := os.Stat(filepath.Join(cmd.Dir, "marker")); !os.IsNotExist(err) {
					t.Fatalf("unexpected marker: %v", err)
				}
			})
		}
	}
}

func TestEvalCommandSubstitution(t *testing.T) {
	shell := testShell(t)
	cases := []struct {
		name, invocation string
		injected         bool
	}{
		{"quoted", `eval "$(printf '%s' "$OUTPUT")"`, false},
		{"unquoted control", `eval $(printf '%s' "$OUTPUT")`, true},
	}
	for name, export := range map[string]ejson2env.ExportFunction{
		"ExportEnv": ejson2env.ExportEnv, "ExportQuiet": ejson2env.ExportQuiet,
	} {
		for _, tc := range cases {
			t.Run(name+"/"+tc.name, func(t *testing.T) {
				dir := t.TempDir()
				// Unquoted substitution expands VALUE='*' to this filename;
				// eval then runs the command in the filename instead of assigning '*'.
				filename := "VALUE=''; : > marker; # '"
				if err := os.WriteFile(filepath.Join(dir, filename), nil, 0600); err != nil {
					t.Fatal(err)
				}
				var output bytes.Buffer
				export(&output, map[string]string{"VALUE": "*"})
				cmd := exec.Command(shell, "-c", tc.invocation+"\nprintf '%s' \"$VALUE\"")
				cmd.Dir = dir
				cmd.Env = []string{"OUTPUT=" + output.String(), "LC_ALL=C"}
				got, err := cmd.CombinedOutput()
				if !tc.injected && (err != nil || string(got) != "*") {
					t.Errorf("quoted eval: got %q, error %v; want literal '*'", got, err)
				}
				// The unquoted control may exit non-zero after running the command.
				_, err = os.Stat(filepath.Join(dir, "marker"))
				if err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
				if injected := err == nil; injected != tc.injected {
					t.Errorf("command executed = %v, want %v", injected, tc.injected)
				}
			})
		}
	}
}
