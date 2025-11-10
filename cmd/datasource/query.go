package datasource

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/internal/grafana"
	"github.com/ynqa/grfnctl/internal/grafana/dashboard"
	"github.com/ynqa/grfnctl/internal/grafana/querymetrics"
	"github.com/ynqa/grfnctl/internal/util"
)

var (
	from     string
	to       string
	interval int64
)

func init() {
	QueryCmd.Flags().StringVar(&from, "from", "now-1h", "Query start time")
	QueryCmd.Flags().StringVar(&to, "to", "now", "Query end time")
	QueryCmd.Flags().Int64Var(&interval, "interval", 5000, "Query interval in milliseconds")
}

var QueryCmd = &cobra.Command{
	Use:   "query [datasource-uid] [query]",
	Short: "Query a Grafana data source",
	Long:  "Debugging tool to validate data source connectivity and query execution in Grafana.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := grafana.NewClientFromEnv()
		if err != nil {
			return err
		}

		datasourceUID := args[0]
		queryStr := args[1]

		payload, err := querymetrics.QueryMetrics(
			client,
			querymetrics.GenerateMetricRequest(
				from,
				to,
				querymetrics.Queries{
					querymetrics.Query{
						RefID: "A",
						Datasource: dashboard.DataSource{
							UID: datasourceUID,
						},
						IntervalMs: interval,
						Expr:       queryStr,
					},
				},
			),
		)
		if err != nil {
			return err
		}

		if err := util.PrintAsJson(payload); err != nil {
			return fmt.Errorf("failed to marshal query metrics payload: %w", err)
		}

		return nil
	},
}
