package util

import (
	"encoding/json"
	"fmt"
)

// PrintAsJson prints the given value as a formatted JSON string.
func PrintAsJson(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data)) // Ensure complete output without truncation
	return nil
}
