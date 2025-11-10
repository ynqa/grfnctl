package dashboard

import (
	"encoding/json"

	"github.com/grafana/grafana-openapi-client-go/models"
)

// Dashboard represents (partially) Grafana dashboard structure.
type Dashboard struct {
	Panels     []Panel         `json:"panels"`
	Templating Templating      `json:"templating"`
	Time       *TimeRange      `json:"time,omitempty"`
	Additional json.RawMessage `json:"-"`
}

// Panel represents (partially) Grafana dashboard panel structure.
type Panel struct {
	Type         PanelType       `json:"type"`
	Targets      []*Target       `json:"targets,omitempty"`
	SnapshotData []*SnapshotData `json:"snapshotData,omitempty"`
	Additional   json.RawMessage `json:"-"`
}

// PanelType represents the type of Grafana dashboard panel.
type PanelType string

const (
	PanelTypeRow PanelType = "row"
)

// Target represents (partially) Grafana dashboard panel target structure.
type Target struct {
	Datasource DataSource      `json:"datasource"`
	Expr       string          `json:"expr,omitempty"`
	RefID      string          `json:"refId,omitempty"`
	Additional json.RawMessage `json:"-"`
}

// SnapshotData captures the schema and data for a Grafana snapshot.
type SnapshotData struct {
	Fields []SnapshotDataField `json:"fields"`
	Meta   *models.FrameMeta   `json:"meta,omitempty"`
	RefID  string              `json:"refId,omitempty"`
}

// SnapshotDataField represents a single field in the snapshot data.
type SnapshotDataField struct {
	Name   string              `json:"name"`
	Type   string              `json:"type,omitempty"`
	Config *models.FieldConfig `json:"config,omitempty"`
	Labels models.FrameLabels  `json:"labels,omitempty"`
	Values []float64           `json:"values,omitempty"`
}

// DataSource represents Grafana dashboard data source structure.
type DataSource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

// Templating represents (partially) Grafana snapshot templating structure.
type Templating struct {
	List []TemplateVariable `json:"list"`
}

// TemplateVariable represents a templated variable definition.
type TemplateVariable struct {
	Name    string                  `json:"name"`
	Current TemplateVariableCurrent `json:"current"`
}

// TemplateVariableCurrent represents current values for a templated variable.
type TemplateVariableCurrent struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

// TimeRange represents (partially) Grafana dashboard time range structure.
type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}
