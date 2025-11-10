package cmd

import (
	"github.com/spf13/cobra"
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Provide Grafana snapshots related commands",
	Long:  "Commands to manage Grafana snapshots, including creation and listing.",
}
