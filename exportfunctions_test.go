package ejson2env_test

import (
	"bytes"
	"testing"

	"github.com/Shopify/ejson2env/v2"
)

func TestExportEnv(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		env      map[string]string
		expected string
	}{
		"empty": {
			env:      map[string]string{},
			expected: "",
		},
		"single key": {
			env: map[string]string{
				"key": "value",
			},
			expected: "export key=value\n",
		},
		"multiple keys": {
			env: map[string]string{
				"key1": "value 1",
				"key2": "value ' with some \" quotes and emoji 🐈",
			},
			expected: "export key1='value 1'\nexport key2='value '\\'' with some \" quotes and emoji 🐈'\n",
		},
		"attempt command injection in key": {
			env: map[string]string{
				"key; touch pwned.txt": "value",
			},
			expected: "",
		},
		"attempt command injection in value": {
			env: map[string]string{
				"key": "value; touch pwned.txt",
			},
			expected: "export key='value; touch pwned.txt'\n",
		},
		"attempt command injection via control characters": {
			env: map[string]string{
				"key": "\bvalue; touch pwned.txt",
			},
			expected: "export key='value; touch pwned.txt'\n",
		},
	}

	for label, tc := range cases {
		tc := tc
		t.Run(label, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			ejson2env.ExportEnv(&buf, tc.env)
			t.Log(buf.String())

			if buf.String() != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, buf.String())
			}
		})
	}
}
