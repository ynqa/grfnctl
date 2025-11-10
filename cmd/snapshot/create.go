package snapshot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/internal/grafana"
	"github.com/ynqa/grfnctl/internal/grafana/dashboard"
	"github.com/ynqa/grfnctl/internal/util"
)

var (
	uid    string
	from   string
	to     string
	expire int64
	vars   []string
	dryRun bool
)

func init() {
	CreateCmd.Flags().StringVarP(&uid, "uid", "u", "", "Dashboard UID to create snapshot from")
	CreateCmd.RegisterFlagCompletionFunc(
		"uid",
		util.DashboardCompletionFunc(),
	)
	CreateCmd.Flags().StringVar(&from, "from", "now-1h", "Snapshot start time")
	CreateCmd.Flags().StringVar(&to, "to", "now", "Snapshot end time")
	CreateCmd.Flags().Int64Var(&expire, "expire", 0, "Snapshot expiration time in seconds (optional)")
	CreateCmd.Flags().StringArrayVar(&vars, "var", nil, "Repeatable key=value pairs for template variables")
	CreateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "If set, the snapshot will not be created, but the converted snapshot JSON will be printed to stdout")
}

var CreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a Grafana snapshot from a dashboard UID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := grafana.NewClientFromEnv()
		if err != nil {
			return err
		}

		boardRaw, err := grafana.GetDashboardBy(client, uid)
		if err != nil {
			return err
		}

		rawJSON, err := json.Marshal(boardRaw.Payload.Dashboard)
		if err != nil {
			return fmt.Errorf("failed to marshal dashboard payload to JSON: %w", err)
		}

		board, err := dashboard.LoadDashboardFrom(rawJSON)
		if err != nil {
			return err
		}

		varsMap := board.TemplatesToKVs()
		// Override with CLI vars.
		for _, pair := range vars {
			if pair == "" {
				continue
			}
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return fmt.Errorf("invalid var format: %q", pair)
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if key == "" {
				return fmt.Errorf("var key cannot be empty in pair %q", pair)
			}
			// Ensures the most recent value wins when keys repeat.
			varsMap[key] = value
		}

		err = grafana.ConvertToSnapshotJSON(
			client,
			board,
			from,
			to,
			varsMap,
		)
		if err != nil {
			return err
		}

		if dryRun {
			if err := util.PrintAsJson(board); err != nil {
				return fmt.Errorf("failed to print snapshot JSON: %w", err)
			}
		} else {
			name := args[0]
			snapshot, err := grafana.CreateSnapshot(
				client,
				&models.CreateDashboardSnapshotCommand{
					Name:      name,
					Dashboard: board,
					Expires:   expire,
				},
			)
			if err != nil {
				return err
			}

			fmt.Println("Snapshot created successfully!")
			fmt.Printf("Snapshot URL: %s\n", snapshot.Payload.URL)
		}
		return nil
	},
}
