package dashboard

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/internal/grafana"
	"github.com/ynqa/grfnctl/internal/tabler"
	"github.com/ynqa/grfnctl/internal/util"
)

var (
	output util.Output = util.OutputTable
)

func init() {
	ListCmd.Flags().VarP(&output, "output", "o", "Output format: json or table")
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Grafana dashboards",
	Long:  "Retrieve and display a list of Grafana dashboards from the configured Grafana instance.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := grafana.NewClientFromEnv()
		if err != nil {
			return err
		}

		searchResult, err := grafana.ListDashboards(client, nil)
		if err != nil {
			return err
		}

		boards := grafana.FilterDashboardsByTypeDB(searchResult)

		switch output {
		case util.OutputJSON:
			if err := util.PrintAsJson(boards); err != nil {
				return fmt.Errorf("failed to marshal dashboards: %w", err)
			}
		case util.OutputTable:
			rows := make([][]string, 0, len(boards))
			for _, board := range boards {
				rows = append(rows, []string{
					board.Title,
					board.UID,
					board.FolderTitle,
					fmt.Sprintf("%s%s", os.Getenv("GRAFANA_SERVER"), board.URL),
				})
			}

			if err := tabler.PrintAsTable(
				tabler.WithHeader([]string{"Title", "UID", "Folder", "URL"}),
				tabler.WithRows(rows),
			); err != nil {
				return fmt.Errorf("failed to print dashboards as table: %w", err)
			}
		}

		return nil
	},
}
