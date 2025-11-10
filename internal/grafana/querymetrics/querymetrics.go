package querymetrics

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/ds"
	"github.com/grafana/grafana-openapi-client-go/models"
)

// QueryMetrics queries metrics from a Grafana data source.
func QueryMetrics(
	client *client.GrafanaHTTPAPI,
	body *models.MetricRequest,
) (*QueryMetricsPayload, error) {
	params := ds.NewQueryMetricsWithExpressionsParams().WithBody(body)
	op := &runtime.ClientOperation{
		ID:                 "queryMetricsWithExpressions",
		Method:             "POST",
		PathPattern:        "/ds/query",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &queryMetricsRawReader{},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}

	result, err := client.Transport.Submit(op)
	if err != nil {
		return nil, err
	}

	metrics, ok := result.(*QueryMetricsPayload)
	if !ok {
		return nil, fmt.Errorf("unexpected response type %T", result)
	}

	return metrics, nil
}

type queryMetricsRawReader struct{}

func (r *queryMetricsRawReader) ReadResponse(
	response runtime.ClientResponse,
	consumer runtime.Consumer,
) (interface{}, error) {
	switch response.Code() {
	case 200, 207:
		var payload QueryMetricsPayload
		if err := consumer.Consume(response.Body(), &payload); err != nil && err != io.EOF {
			return nil, err
		}
		return &payload, nil
	default:
		data, err := io.ReadAll(response.Body())
		if err != nil {
			return nil, err
		}
		msg := strings.TrimSpace(string(data))
		if msg == "" {
			msg = fmt.Sprintf("status %d", response.Code())
		} else {
			msg = fmt.Sprintf("status %d: %s", response.Code(), msg)
		}
		return nil, fmt.Errorf("query metrics failed: %s", msg)
	}
}
