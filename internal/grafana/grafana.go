package grafana

import (
	"github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/dashboards"
	"github.com/grafana/grafana-openapi-client-go/client/datasources"
	"github.com/grafana/grafana-openapi-client-go/client/folders"
	"github.com/grafana/grafana-openapi-client-go/client/search"
	"github.com/grafana/grafana-openapi-client-go/client/signed_in_user"
	"github.com/grafana/grafana-openapi-client-go/client/snapshots"
	"github.com/grafana/grafana-openapi-client-go/models"
)

// NOTE: The Dashboard API provided by grafana/grafana-openapi-client-go
// (version v0.0.0-20250925215610-d92957c70d5c) includes an outdated endpoint.
// - https://grafana.com/docs/grafana/v11.6/developers/http_api/dashboard
// Please consider replacing it with the new endpoint when it becomes available.
// - https://grafana.com/docs/grafana/v12.2/developers/http_api/dashboard

// GetDashboard retrieves a Grafana dashboard by its UID.
func GetDashboardBy(
	client *client.GrafanaHTTPAPI,
	uid string,
) (*dashboards.GetDashboardByUIDOK, error) {
	return client.Dashboards.GetDashboardByUID(uid)
}

type DashboardType string

const (
	DashboardTypeDB     DashboardType = "dash-db"
	DashboardTypeFolder DashboardType = "dash-folder"
)

// ListDashboards retrieves a list of Grafana dashboards based on the provided parameters.
func ListDashboards(
	client *client.GrafanaHTTPAPI,
	params *search.SearchParams,
) (*search.SearchOK, error) {
	return client.Search.Search(params)
}

// FilterDashboardsByTypeDB returns dashboards filtered by "dash-db" type.
func FilterDashboardsByTypeDB(boards *search.SearchOK) []*models.Hit {
	if boards == nil || boards.Payload == nil {
		return nil
	}

	filtered := make([]*models.Hit, 0, len(boards.Payload))
	for _, board := range boards.Payload {
		if board == nil {
			continue
		}

		if DashboardType(board.Type) == DashboardTypeDB {
			filtered = append(filtered, board)
		}
	}

	return filtered
}

// CreateOrUpdateDashboard creates or updates a Grafana dashboard.
func CreateOrUpdateDashboard(
	client *client.GrafanaHTTPAPI,
	body *models.SaveDashboardCommand,
) (*dashboards.PostDashboardOK, error) {
	return client.Dashboards.PostDashboard(body)
}

// ListDataSources retrieves a list of Grafana data sources.
func ListDataSources(
	client *client.GrafanaHTTPAPI,
) (*datasources.GetDataSourcesOK, error) {
	return client.Datasources.GetDataSources()
}

// ListFolders retrieves a list of Grafana folders based on the provided parameters.
func ListFolders(
	client *client.GrafanaHTTPAPI,
	params *folders.GetFoldersParams,
) (*folders.GetFoldersOK, error) {
	return client.Folders.GetFolders(params)
}

// CreateSnapshot creates a new Grafana dashboard snapshot.
func CreateSnapshot(
	client *client.GrafanaHTTPAPI,
	body *models.CreateDashboardSnapshotCommand,
) (*snapshots.CreateDashboardSnapshotOK, error) {
	return client.Snapshots.CreateDashboardSnapshot(body)
}

// ListSnapshots retrieves a list of Grafana snapshots based on the provided parameters.
// Note: `Admin` privileges are required to list Grafana snapshots.
func ListSnapshots(
	client *client.GrafanaHTTPAPI,
	params *snapshots.SearchDashboardSnapshotsParams,
) (*snapshots.SearchDashboardSnapshotsOK, error) {
	return client.Snapshots.SearchDashboardSnapshots(params)
}

// DeleteSnapshotBy deletes a Grafana snapshot by its key (UID).
func DeleteSnapshotBy(
	client *client.GrafanaHTTPAPI,
	key string,
) (*snapshots.DeleteDashboardSnapshotOK, error) {
	return client.Snapshots.DeleteDashboardSnapshot(key)
}

// GetCurrentUser retrieves information about the currently authenticated user.
func GetCurrentUser(
	client *client.GrafanaHTTPAPI,
) (*signed_in_user.GetSignedInUserOK, error) {
	return client.SignedInUser.GetSignedInUser()
}
