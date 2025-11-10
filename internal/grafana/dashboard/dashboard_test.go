package dashboard

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseTargetFromJSON(t *testing.T) {
	inputJSON := []byte(`{
		"datasource": {
			"type": "prometheus",
			"uid": "prom"
		},
		"expr": "rate(requests_total[5m])",
		"refId": "A",
		"legendFormat": "{{instance}}",
		"interval": "15s"
	}`)

	var target Target
	if err := json.Unmarshal(inputJSON, &target); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	wantDatasource := DataSource{
		Type: "prometheus",
		UID:  "prom",
	}
	if target.Datasource != wantDatasource {
		t.Errorf("Datasource mismatch.\nWant: %v\nGot:  %v", wantDatasource, target.Datasource)
	}

	if target.Expr != "rate(requests_total[5m])" {
		t.Errorf("Expr mismatch.\nWant: %v\nGot:  %v", "rate(requests_total[5m])", target.Expr)
	}

	if target.RefID != "A" {
		t.Errorf("RefID mismatch.\nWant: %v\nGot:  %v", "A", target.RefID)
	}
}

func TestLoadDashboardFromPreservesAdditional(t *testing.T) {
	inputJSON := []byte(`{
		"panels": [
			{
				"type": "row",
				"targets": [
					{
						"datasource": {
							"type": "prometheus",
							"uid": "prom"
						},
						"expr": "rate(requests_total[5m])",
						"legendFormat": "{{instance}}"
					}
				],
				"title": "Row Panel",
				"gridPos": {
					"h": 1,
					"w": 24,
					"x": 0,
					"y": 0
				}
			}
		],
		"templating": {
			"list": [
				{
					"name": "env",
					"current": {
						"text": "prod",
						"value": "prod"
					}
				}
			]
		},
		"schemaVersion": 37,
		"refresh": "5s"
	}`)

	board, err := LoadDashboardFrom(inputJSON)
	if err != nil {
		t.Fatalf("LoadDashboardFrom returned error: %v", err)
	}

	gotJSON, err := json.Marshal(board)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	var want map[string]any
	if err := json.Unmarshal(inputJSON, &want); err != nil {
		t.Fatalf("json.Unmarshal on input returned error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(gotJSON, &got); err != nil {
		t.Fatalf("json.Unmarshal on output returned error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Dashboard JSON mismatch.\nWant: %v\nGot:  %v", want, got)
	}
}

func TestTemplatesToKVs(t *testing.T) {
	inputJSON := []byte(`{
		"templating": {
			"list": [
				{
					"name": "env",
					"current": {
						"text": "prod",
						"value": "prod"
					}
				},
				{
					"name": "region",
					"current": {
						"text": "us-west",
						"value": "us-west"
					}
				}
			]
		}
	}`)

	board, err := LoadDashboardFrom(inputJSON)
	if err != nil {
		t.Fatalf("LoadDashboardFrom returned error: %v", err)
	}

	varsMap := board.TemplatesToKVs()
	expected := map[string]string{
		"env":    "prod",
		"region": "us-west",
	}

	if !reflect.DeepEqual(varsMap, expected) {
		t.Errorf("TemplatesToKVs result mismatch.\nWant: %v\nGot:  %v", expected, varsMap)
	}
}
