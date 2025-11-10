package dashboard

import (
	"encoding/json"
)

type dashboardPayload struct {
	Panels     []Panel    `json:"panels"`
	Templating Templating `json:"templating"`
	Time       *TimeRange `json:"time,omitempty"`
}

// UnmarshalJSON binds known fields and stores the remaining payload in Additional.
func (d *Dashboard) UnmarshalJSON(data []byte) error {
	var aux dashboardPayload
	if err := decodeWithAdditional(
		data,
		&aux,
		&d.Additional,
		"panels",
		"templating",
		"time",
	); err != nil {
		return err
	}

	d.Panels = aux.Panels
	d.Templating = aux.Templating
	d.Time = aux.Time
	return nil
}

// MarshalJSON merges known fields and Additional into a single JSON object.
func (d Dashboard) MarshalJSON() ([]byte, error) {
	base, err := json.Marshal(dashboardPayload{
		Panels:     d.Panels,
		Templating: d.Templating,
		Time:       d.Time,
	})
	if err != nil {
		return nil, err
	}

	return mergeWithAdditional(base, d.Additional)
}

type panelPayload struct {
	Type         PanelType       `json:"type"`
	Targets      []*Target       `json:"targets"`
	SnapshotData []*SnapshotData `json:"snapshotData,omitempty"`
}

// UnmarshalJSON binds known fields and stores the remaining payload in Additional.
func (p *Panel) UnmarshalJSON(data []byte) error {
	var aux panelPayload
	if err := decodeWithAdditional(
		data,
		&aux,
		&p.Additional,
		"type",
		"targets",
		"snapshotData",
	); err != nil {
		return err
	}

	p.Type = aux.Type
	p.Targets = aux.Targets
	p.SnapshotData = aux.SnapshotData
	return nil
}

// MarshalJSON merges known fields and Additional into a single JSON object.
func (p Panel) MarshalJSON() ([]byte, error) {
	base, err := json.Marshal(panelPayload{
		Type:         p.Type,
		Targets:      p.Targets,
		SnapshotData: p.SnapshotData,
	})
	if err != nil {
		return nil, err
	}

	return mergeWithAdditional(base, p.Additional)
}

type targetPayload struct {
	Datasource DataSource `json:"datasource"`
	Expr       string     `json:"expr"`
	RefID      string     `json:"refId,omitempty"`
}

// UnmarshalJSON binds known fields and stores the remaining payload in Additional.
func (t *Target) UnmarshalJSON(data []byte) error {
	var aux targetPayload
	if err := decodeWithAdditional(
		data,
		&aux,
		&t.Additional,
		"datasource",
		"expr",
		"refId",
	); err != nil {
		return err
	}

	t.Datasource = aux.Datasource
	t.Expr = aux.Expr
	t.RefID = aux.RefID
	return nil
}

// MarshalJSON merges known fields and Additional into a single JSON object.
func (t Target) MarshalJSON() ([]byte, error) {
	base, err := json.Marshal(targetPayload{
		Datasource: t.Datasource,
		Expr:       t.Expr,
		RefID:      t.RefID,
	})
	if err != nil {
		return nil, err
	}

	return mergeWithAdditional(base, t.Additional)
}

// mergeWithAdditional combines base JSON with additional fields.
func mergeWithAdditional(base []byte, additional json.RawMessage) ([]byte, error) {
	if len(additional) == 0 {
		return base, nil
	}

	var baseMap map[string]json.RawMessage
	if err := json.Unmarshal(base, &baseMap); err != nil {
		return nil, err
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(additional, &rawMap); err != nil {
		return nil, err
	}

	for k, v := range rawMap {
		baseMap[k] = v
	}

	return json.Marshal(baseMap)
}

// decodeWithAdditional binds known fields and captures the remaining payload into additional.
func decodeWithAdditional(data []byte, dest any, additional *json.RawMessage, knownKeys ...string) error {
	if err := json.Unmarshal(data, dest); err != nil {
		return err
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	for _, key := range knownKeys {
		delete(rawMap, key)
	}

	if len(rawMap) == 0 {
		*additional = (*additional)[:0]
		return nil
	}

	extra, err := json.Marshal(rawMap)
	if err != nil {
		return err
	}

	*additional = append((*additional)[:0], extra...)
	return nil
}
