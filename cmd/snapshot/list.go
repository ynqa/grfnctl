package snapshot

import (
	"fmt"

	"github.com/grafana/grafana-openapi-client-go/client/snapshots"
	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/internal/grafana"
	"github.com/ynqa/grfnctl/internal/tabler"
	"github.com/ynqa/grfnctl/internal/util"
)

var (
	limit  int64
	query  string
	output util.Output = util.OutputTable
)

func init() {
	ListCmd.Flags().Int64VarP(&limit, "limit", "l", 1000, "Limit the number of snapshots to retrieve")
	ListCmd.Flags().StringVarP(&query, "query", "q", "", "Query string to filter snapshots by name")
	ListCmd.Flags().VarP(&output, "output", "o", "Output format: json or table")
}

// ListCmd represents the `snapshot list` command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Grafana snapshots",
	Long:  "Retrieve and display a list of Grafana snapshots from the configured Grafana instance.",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := snapshots.SearchDashboardSnapshotsParams{
			Limit: &limit,
			Query: &query,
		}

		client, err := grafana.NewClientFromEnv()
		if err != nil {
			return err
		}

		snaps, err := grafana.ListSnapshots(client, &params)
		if err != nil {
			return err
		}

		switch output {
		case util.OutputJSON:
			if err := util.PrintAsJson(snaps.Payload); err != nil {
				return fmt.Errorf("failed to marshal snapshots: %w", err)
			}
		case util.OutputTable:
			rows := make([][]string, 0, len(snaps.Payload))
			for _, snap := range snaps.Payload {
				rows = append(rows, []string{
					snap.Name,
					snap.Key,
					util.FormatDateTimeInLocal(snap.Created),
					util.FormatDateTimeInLocal(snap.Expires),
				})
			}

			if err := tabler.PrintAsTable(
				tabler.WithHeader([]string{"Name", "Key", "Created", "Expires"}),
				tabler.WithRows(rows),
			); err != nil {
				return fmt.Errorf("failed to render table: %w", err)
			}
		}

		return nil
	},
}
