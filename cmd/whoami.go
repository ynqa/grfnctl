package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/internal/grafana"
	"github.com/ynqa/grfnctl/internal/tabler"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display the current user information",
	Long:  "This command retrieves and displays information about the currently authenticated user.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := grafana.NewClientFromEnv()
		if err != nil {
			return err
		}

		user, err := grafana.GetCurrentUser(client)
		if err != nil {
			return err
		}

		rows := make([][]string, 0, 3)
		rows = append(rows, []string{"ID", user.Payload.Login})
		rows = append(rows, []string{"Name", user.Payload.Name})
		rows = append(rows, []string{"UID", user.Payload.UID})

		if err := tabler.PrintAsTable(
			tabler.WithHeader([]string{"Attribute", "Value"}),
			tabler.WithRows(rows),
		); err != nil {
			return fmt.Errorf("failed to print user info as table: %w", err)
		}
		return nil
	},
}
