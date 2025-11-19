package ejson2env_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/Shopify/ejson2env/v2"
)

func TestExport(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		env           map[string]string
		expected      string
		expectedQuiet string
	}{
		"empty": {
			env:      map[string]string{},
			expected: "",
		},
		"single key": {
			env: map[string]string{
				"key": "value",
			},
			expected:      "export key=value\n",
			expectedQuiet: "key=value\n",
		},
		"attempt command injection in key": {
			env: map[string]string{
				"key; touch pwned.txt": "value",
			},
			expected:      "",
			expectedQuiet: "",
		},
		"newline in key": {
			env: map[string]string{
				"touch pwned.txt;\ndummy": "value",
			},
			expected:      "",
			expectedQuiet: "",
		},
		"attempt command injection in value": {
			env: map[string]string{
				"key": "value; touch pwned.txt",
			},
			expected:      "export key='value; touch pwned.txt'\n",
			expectedQuiet: "key='value; touch pwned.txt'\n",
		},
		"attempt command injection via control characters": {
			env: map[string]string{
				"key": "\bvalue; touch pwned.txt",
			},
			expected:      "export key='value; touch pwned.txt'\n",
			expectedQuiet: "key='value; touch pwned.txt'\n",
		},
		"newline in value": {
			env: map[string]string{
				"key": "value\nnewline",
			},
			expected:      "export key='value\nnewline'\n",
			expectedQuiet: "key='value\nnewline'\n",
		},
		"escaped newlines in value": {
			env: map[string]string{
				"key": "value\\nnewline",
			},
			expected:      "export key='value\\nnewline'\n",
			expectedQuiet: "key='value\\nnewline'\n",
		},
	}

	for label, tc := range cases {
		tc := tc
		t.Run(label, func(t *testing.T) {
			t.Parallel()

			t.Run("ExportEnv", func(t *testing.T) {
				var buf bytes.Buffer
				ejson2env.ExportEnv(&buf, tc.env)
				t.Log(buf.String())

				if buf.String() != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, buf.String())
				}
			})

			t.Run("ExportQuiet", func(t *testing.T) {
				var buf bytes.Buffer
				ejson2env.ExportQuiet(&buf, tc.env)
				t.Log(buf.String())

				if buf.String() != tc.expectedQuiet {
					t.Errorf("expected %q, got %q", tc.expectedQuiet, buf.String())
				}
			})
		})
	}
}

// TestExportEnvMultipleKeys tests exporting multiple environment variables
// with map key order independence
func TestExportEnvMultipleKeys(t *testing.T) {
	fmt.Println("===== RUNNING TestExportEnvMultipleKeys TEST =====")
	t.Parallel()

	env := map[string]string{
		"key1": "value 1",
		"key2": "value ' with some \" quotes and emoji 🐈",
	}

	var buf bytes.Buffer
	ejson2env.ExportEnv(&buf, env)
	output := buf.String()
	t.Log(output)

	// Check for each expected line individually
	if !strings.Contains(output, "export key1='value 1'") {
		t.Errorf("output missing 'export key1='value 1''")
	}
	if !strings.Contains(output, "export key2='value '\"'\"' with some \" quotes and emoji 🐈'") {
		t.Errorf("output missing key2 with proper escaping")
	}
}

// TestGitHubActionsMaskWrapper tests that the GitHub Actions wrapper
// outputs mask commands for non-underscore-prefixed keys
func TestGitHubActionsMaskWrapper(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		env         map[string]string
		expected    []string // lines we expect to see
		notExpected []string // lines we should NOT see
	}{
		"masks secret values": {
			env: map[string]string{
				"SECRET_KEY": "my-secret-value",
				"API_TOKEN":  "token123",
			},
			expected: []string{
				"echo \"::add-mask::my-secret-value\"",
				"echo \"::add-mask::token123\"",
				"export API_TOKEN=token123",
				"export SECRET_KEY=my-secret-value",
			},
		},
		"does not mask underscore-prefixed keys": {
			env: map[string]string{
				"_PUBLIC_KEY": "public-value",
				"SECRET_KEY":  "secret-value",
			},
			expected: []string{
				"echo \"::add-mask::secret-value\"",
				"export SECRET_KEY=secret-value",
				"export _PUBLIC_KEY=public-value",
			},
			notExpected: []string{
				"echo \"::add-mask::public-value\"",
			},
		},
		"handles special characters in values": {
			env: map[string]string{
				"SECRET": "value with spaces",
			},
			expected: []string{
				"echo \"::add-mask::'value with spaces'\"",
				"export SECRET='value with spaces'",
			},
		},
		"empty map": {
			env:      map[string]string{},
			expected: []string{},
		},
	}

	for label, tc := range cases {
		tc := tc
		t.Run(label, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			wrapped := ejson2env.GitHubActionsMaskWrapper(ejson2env.ExportEnv)
			wrapped(&buf, tc.env)
			output := buf.String()
			t.Log(output)

			// Check for expected lines
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("output missing expected line: %q", expected)
				}
			}

			// Check for lines that should NOT be present
			for _, notExpected := range tc.notExpected {
				if strings.Contains(output, notExpected) {
					t.Errorf("output contains unexpected line: %q", notExpected)
				}
			}
		})
	}
}

// TestGitHubActionsMaskWithTrimUnderscore tests that the wrappers work correctly
// when chained together. The mask wrapper should be applied to the base function first,
// then the trim wrapper wraps that, matching the order in main.go
func TestGitHubActionsMaskWithTrimUnderscore(t *testing.T) {
	t.Parallel()

	env := map[string]string{
		"_PUBLIC_KEY": "public-value",
		"SECRET_KEY":  "secret-value",
	}

	var buf bytes.Buffer
	// Match the order from main.go:
	// 1. Start with base export function
	// 2. Apply trim wrapper (inner)
	// 3. Apply mask wrapper (outer)
	exportFunc := ejson2env.ExportEnv
	exportFunc = ejson2env.TrimLeadingUnderscoreExportWrapper(exportFunc)
	exportFunc = ejson2env.GitHubActionsMaskWrapper(exportFunc)

	exportFunc(&buf, env)
	output := buf.String()
	t.Log(output)

	// Should mask the secret but not the public value
	if !strings.Contains(output, "echo \"::add-mask::secret-value\"") {
		t.Errorf("output missing mask command for secret-value")
	}

	if strings.Contains(output, "echo \"::add-mask::public-value\"") {
		t.Errorf("output should not mask public-value (underscore-prefixed)")
	}

	// Both keys should be exported with underscore trimmed
	if !strings.Contains(output, "export PUBLIC_KEY=public-value") {
		t.Errorf("output missing export with trimmed underscore for PUBLIC_KEY")
	}

	if !strings.Contains(output, "export SECRET_KEY=secret-value") {
		t.Errorf("output missing export for SECRET_KEY")
	}
}
