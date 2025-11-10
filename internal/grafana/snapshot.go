package grafana

import (
	"errors"

	"github.com/grafana/grafana-openapi-client-go/client"

	"github.com/ynqa/grfnctl/internal/grafana/dashboard"
	"github.com/ynqa/grfnctl/internal/grafana/querymetrics"
	"github.com/ynqa/grfnctl/internal/grafana/render"
)

// ConvertToSnapshotJSON converts the given dashboard to use snapshot data by
// replacing templated variables and executing metric queries.
func ConvertToSnapshotJSON(
	client *client.GrafanaHTTPAPI,
	board *dashboard.Dashboard,
	from string,
	to string,
	kvs map[string]string,
) error {
	// First, replace templated variables and execute queries.
	for i := range board.Panels {
		panel := &board.Panels[i]
		// Skip rows panels.
		if panel.Type == dashboard.PanelTypeRow {
			continue
		}

		snapshotData := []*dashboard.SnapshotData{}

		// Replace templated variables in targets.
		for _, target := range panel.Targets {
			// Process only supported data source types.
			if target.Datasource.Type != "prometheus" {
				return errors.New("only `prometheus` data source is supported in snapshot")
			}

			// Replace variables in expr and datasource UID.
			datasourceUID := render.RenderVariables(target.Datasource.UID, kvs)
			expr := render.RenderVariables(target.Expr, kvs)
			refID := target.RefID

			// Execute query to get metric data.
			payload, err := querymetrics.QueryMetrics(
				client,
				querymetrics.GenerateMetricRequest(
					from,
					to,
					querymetrics.Queries{
						querymetrics.Query{
							RefID: refID,
							Datasource: dashboard.DataSource{
								UID: datasourceUID,
							},
							IntervalMs: 5000,
							Expr:       expr,
						},
					},
				),
			)
			if err != nil {
				return err
			}

			for _, frame := range payload.Results[refID].Frames {
				// NOTE:
				// - frame.Data.Values[0] means timestamp
				// - frame.Data.Values[1] means values for query results
				if len(frame.Data.Values) != 0 {
					snapshotData = append(snapshotData, &dashboard.SnapshotData{
						Meta:  frame.Schema.Meta,
						RefID: frame.Schema.RefID,
						Fields: func() []dashboard.SnapshotDataField {
							fields := []dashboard.SnapshotDataField{}
							for _, field := range frame.Schema.Fields {
								fields = append(fields, dashboard.SnapshotDataField{
									Name:   "Time",
									Type:   "time",
									Config: field.Config,
									Labels: field.Labels,
									Values: frame.Data.Values[0],
								})
								fields = append(fields, dashboard.SnapshotDataField{
									Name:   field.Name,
									Type:   field.Type,
									Config: field.Config,
									Labels: field.Labels,
									Values: frame.Data.Values[1],
								})
							}
							return fields
						}(),
					})
				}
			}
		}

		panel.Targets = nil
		panel.SnapshotData = snapshotData
	}

	// Update dashboard time range.
	board.Time = &dashboard.TimeRange{
		From: from,
		To:   to,
	}

	return nil
}
