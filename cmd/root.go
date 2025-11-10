package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/cmd/dashboard"
	"github.com/ynqa/grfnctl/cmd/datasource"
	"github.com/ynqa/grfnctl/cmd/snapshot"
)

var rootCmd = &cobra.Command{
	Use:   "grfnctl",
	Short: "Grafana CLI tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(dashboardCmd)
	dashboardCmd.AddCommand(dashboard.ListCmd)
	dashboardCmd.AddCommand(dashboard.ExportCmd)
	dashboardCmd.AddCommand(dashboard.ApplyCmd)
	rootCmd.AddCommand(datasourceCmd)
	datasourceCmd.AddCommand(datasource.QueryCmd)
	datasourceCmd.AddCommand(datasource.ListCmd)
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.AddCommand(snapshot.ListCmd)
	snapshotCmd.AddCommand(snapshot.DeleteCmd)
	snapshotCmd.AddCommand(snapshot.CreateCmd)
	rootCmd.AddCommand(whoamiCmd)
}
