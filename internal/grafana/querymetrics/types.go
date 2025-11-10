package querymetrics

import "github.com/grafana/grafana-openapi-client-go/models"

// QueryMetricsPayload mirrors models.QueryDataResponse but fixes the JSON tags
// so they match Grafana's lowercase schema. In
// grafana/grafana-openapi-client-go v0.0.0-20250925215610-d92957c70d5c,
// models.DataResponse (models/data_response.go) and its nested Frame types
// (e.g. models.Frame in models/frame.go) still declare capitalized tags such as
// `json:"Error"` and `json:"Fields"`, which prevents the generated structs from
// unmarshalling real Grafana payloads that use `error`, `frames`, and other
// lowercase keys. These thin wrappers preserve the upstream structure while
// correcting the tags.
type QueryMetricsPayload struct {
	Results QueryMetricsResponses `json:"results,omitempty"`
}

// QueryMetricsResponses maps RefIDs to their corresponding query data.
type QueryMetricsResponses map[string]QueryMetricsDataResponse

// QueryMetricsDataResponse matches the structure of models.DataResponse while
// honoring Grafana's lowercase field names.
type QueryMetricsDataResponse struct {
	Error       string             `json:"error,omitempty"`
	ErrorSource models.Source      `json:"errorSource,omitempty"`
	Frames      QueryMetricsFrames `json:"frames,omitempty"`
	Status      models.Status      `json:"status,omitempty"`
}

// QueryMetricsFrames describes Grafana's schema/data frame representation.
type QueryMetricsFrames []*QueryMetricsFrame

// QueryMetricsFrame splits the Grafana frame into schema and data blocks.
type QueryMetricsFrame struct {
	Data   QueryMetricsFrameData   `json:"data"`
	Schema QueryMetricsFrameSchema `json:"schema"`
}

// QueryMetricsFrameData captures the columnar values emitted by Grafana.
type QueryMetricsFrameData struct {
	Values [][]float64 `json:"values"`
}

// QueryMetricsFrameSchema describes metadata for each Grafana frame.
type QueryMetricsFrameSchema struct {
	Fields []QueryMetricsSchemaField `json:"fields"`
	Meta   *models.FrameMeta         `json:"meta,omitempty"`
	RefID  string                    `json:"refId,omitempty"`
}

// QueryMetricsSchemaField models the per-column metadata inside a frame.
type QueryMetricsSchemaField struct {
	Name     string              `json:"name"`
	Type     string              `json:"type,omitempty"`
	TypeInfo map[string]any      `json:"typeInfo,omitempty"`
	Config   *models.FieldConfig `json:"config,omitempty"`
	Labels   models.FrameLabels  `json:"labels,omitempty"`
}
