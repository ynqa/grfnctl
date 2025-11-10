package dashboard

import (
	"encoding/json"
	"fmt"
)

// LoadDashboardFrom loads the dashboard templating section from the given dashboard payload.
func LoadDashboardFrom(rawJSON []byte) (*Dashboard, error) {
	var d Dashboard

	if err := json.Unmarshal(rawJSON, &d); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dashboard JSON: %w", err)
	}

	return &d, nil
}

// TemplatesToKVs converts the dashboard templating variables to a key-value map.
func (d *Dashboard) TemplatesToKVs() map[string]string {
	kvs := make(map[string]string)

	for _, v := range d.Templating.List {
		kvs[v.Name] = v.Current.Value
	}

	return kvs
}
