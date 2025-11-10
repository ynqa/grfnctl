package datasource

import (
	"fmt"

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
	Short: "List Grafana data sources",
	Long:  "Retrieve and display a list of Grafana data sources from the configured Grafana instance.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := grafana.NewClientFromEnv()
		if err != nil {
			return err
		}

		datasources, err := grafana.ListDataSources(client)
		if err != nil {
			return err
		}

		switch output {
		case util.OutputJSON:
			if err := util.PrintAsJson(datasources.Payload); err != nil {
				return fmt.Errorf("failed to marshal data sources: %w", err)
			}
		case util.OutputTable:
			rows := make([][]string, 0, len(datasources.Payload))
			for _, ds := range datasources.Payload {
				rows = append(rows, []string{
					ds.Name,
					ds.UID,
					ds.Type,
				})
			}

			if err := tabler.PrintAsTable(
				tabler.WithHeader([]string{"Name", "UID", "Type"}),
				tabler.WithRows(rows),
			); err != nil {
				return fmt.Errorf("failed to print data sources as table: %w", err)
			}
		}

		return nil
	},
}
