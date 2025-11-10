package cmd

import (
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Provide Grafana dashboard related commands",
	Long:  "Commands to manage Grafana dashboards, including creation and listing.",
}
