package cmd

import (
	"github.com/spf13/cobra"
)

var datasourceCmd = &cobra.Command{
	Use:   "datasource",
	Short: "Provide Grafana data source related commands",
	Long:  "Commands to manage Grafana data sources, including creation and listing.",
}
