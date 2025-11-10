package render

import (
	"regexp"
)

// Reference:
// - https://github.com/grafana/grafana/blob/v12.2.1/public/app/features/variables/utils.ts#L26-L33
var (
	variableRegex = regexp.MustCompile(`\$(\w+)|\[\[(\w+?)(?::(\w+))?\]\]|\${(\w+)(?:\.([^:^}]+))?(?::([^}]+))?}`)
)

// RenderVariables replaces templated variables in string using the provided KVs.
func RenderVariables(input string, kvs map[string]string) string {
	if len(kvs) == 0 {
		return input
	}

	return variableRegex.ReplaceAllStringFunc(input, func(match string) string {
		submatches := variableRegex.FindStringSubmatch(match)
		if len(submatches) == 0 {
			return match
		}

		// Capture groups 1, 2, and 4 of variableRegex respectively represent
		// $var, [[var]], and ${var:format} syntax.
		varNames := []int{1, 2, 4}
		for _, idx := range varNames {
			if idx < len(submatches) {
				if value, ok := kvs[submatches[idx]]; ok {
					return value
				}
			}
		}

		return match
	})
}
