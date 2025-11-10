package render

import (
	"testing"
)

func TestRenderVariables(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		kvs      map[string]string
		expected string
	}{
		{
			name:     "dollar syntax",
			input:    `namespace="$namespace"`,
			kvs:      map[string]string{"namespace": "prod"},
			expected: `namespace="prod"`,
		},
		{
			name:     "double bracket syntax",
			input:    `[[namespace]] pods`,
			kvs:      map[string]string{"namespace": "staging"},
			expected: `staging pods`,
		},
		{
			name:     "curly syntax with formatting",
			input:    `title ${namespace:raw}`,
			kvs:      map[string]string{"namespace": "dev"},
			expected: `title dev`,
		},
		{
			name:     "multiple variables",
			input:    `namespace="$namespace" cluster="$cluster"`,
			kvs:      map[string]string{"namespace": "prod", "cluster": "east"},
			expected: `namespace="prod" cluster="east"`,
		},
		{
			name:     "unmatched variable remains",
			input:    `cluster="$cluster"`,
			kvs:      map[string]string{"namespace": "prod"},
			expected: `cluster="$cluster"`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := RenderVariables(tc.input, tc.kvs)
			if got != tc.expected {
				t.Fatalf("replaceVariableValue(%q, %v) = %q, want %q", tc.input, tc.kvs, got, tc.expected)
			}
		})
	}
}
