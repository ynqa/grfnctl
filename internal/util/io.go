package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadDashboardFile reads a Grafana dashboard JSON file from the specified path
// and unmarshals it into a dynamic map structure.
func ReadDashboardFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dashboard map[string]interface{}
	if err := json.Unmarshal(data, &dashboard); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dashboard JSON: %w", err)
	}

	return dashboard, nil
}
