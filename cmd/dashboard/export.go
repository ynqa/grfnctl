package dashboard

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/internal/grafana"
	"github.com/ynqa/grfnctl/internal/util"
)

var ExportCmd = &cobra.Command{
	Use:   "export [uid]",
	Short: "Export Grafana dashboards",
	Long:  "Export Grafana dashboards from the configured Grafana instance.",
	// Provide Grafana dashboard UIDs for tab completion.
	ValidArgsFunction: util.DashboardCompletionFunc(util.SkipCompletionWhenArgsProvided()),
	RunE: func(cmd *cobra.Command, args []string) error {
		uid := args[0]

		client, err := grafana.NewClientFromEnv()
		if err != nil {
			return err
		}

		dashboard, err := grafana.GetDashboardBy(client, uid)
		if err != nil {
			return err
		}

		payload := dashboard.GetPayload()
		if payload == nil {
			return fmt.Errorf("dashboard %s not found", uid)
		}

		if err := util.PrintAsJson(payload); err != nil {
			return err
		}

		return nil
	},
}
