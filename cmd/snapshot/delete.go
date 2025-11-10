package snapshot

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/internal/grafana"
	"github.com/ynqa/grfnctl/internal/util"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [uid]",
	Short: "Delete a Grafana snapshot by its UID",
	Args:  cobra.ExactArgs(1),
	// Provide Grafana snapshot Keys for tab completion.
	ValidArgsFunction: util.SnapshotCompletionFunc(util.SkipCompletionWhenArgsProvided()),
	RunE: func(cmd *cobra.Command, args []string) error {
		uid := args[0]

		client, err := grafana.NewClientFromEnv()
		if err != nil {
			return err
		}

		_, err = grafana.DeleteSnapshotBy(client, uid)
		if err != nil {
			return err
		}

		fmt.Printf("Snapshot with UID '%s' has been deleted successfully.\n", uid)
		return nil
	},
}
