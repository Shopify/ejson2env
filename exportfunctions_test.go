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
			expected:      "export key=valuenewline\n",
			expectedQuiet: "key=valuenewline\n",
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
	if !strings.Contains(output, "export key2='value '\\'' with some \" quotes and emoji 🐈'") {
		t.Errorf("output missing key2 with proper escaping")
	}
}
