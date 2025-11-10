package grafana

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-openapi-client-go/client"
)

func NewClientFromEnv() (*client.GrafanaHTTPAPI, error) {
	cfg, err := newConfigFromEnv()
	if err != nil {
		return nil, err
	}

	return client.NewHTTPClientWithConfig(strfmt.Default, cfg), nil
}

func newConfigFromEnv() (*client.TransportConfig, error) {
	server := os.Getenv("GRAFANA_SERVER")
	if server == "" {
		return nil, errors.New("GRAFANA_SERVER environment variable is required")
	}

	grafanaURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	cfg := client.DefaultTransportConfig().
		WithHost(grafanaURL.Host).
		WithBasePath(strings.TrimLeft(grafanaURL.Path+"/api", "/")).
		WithSchemes([]string{grafanaURL.Scheme})

	apiKey := os.Getenv("GRAFANA_TOKEN")
	user := os.Getenv("GRAFANA_USER")
	password := os.Getenv("GRAFANA_PASSWORD")

	if apiKey != "" {
		cfg.APIKey = apiKey
	} else if user != "" && password != "" {
		cfg.BasicAuth = url.UserPassword(user, password)
	}

	if orgIDStr := os.Getenv("GRAFANA_ORG_ID"); orgIDStr != "" {
		orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
		if err != nil {
			return nil, errors.New("invalid GRAFANA_ORG_ID: must be a number")
		}
		cfg.OrgID = orgID
	}

	return cfg, nil
}
