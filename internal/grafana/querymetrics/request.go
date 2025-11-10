package querymetrics

import (
	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/ynqa/grfnctl/internal/grafana/dashboard"
)

// Query expresses the minimal JSON structure required for Grafana metric queries.
type Query struct {
	RefID      string               `json:"refId"`
	Datasource dashboard.DataSource `json:"datasource"`
	IntervalMs int64                `json:"intervalMs"`
	Expr       string               `json:"expr,omitempty"`
}

type Queries []Query

// ToJSONSlice converts the Queries to a slice of models.JSON.
func (q Queries) ToJSONSlice() []models.JSON {
	result := make([]models.JSON, 0, len(q))
	for _, query := range q {
		result = append(result, models.JSON(query))
	}
	return result
}

// GenerateMetricRequest creates a MetricRequest for the given time range and queries.
func GenerateMetricRequest(
	from string,
	to string,
	queries Queries,
) *models.MetricRequest {
	return &models.MetricRequest{
		From:    &from,
		To:      &to,
		Queries: queries.ToJSONSlice(),
	}
}
