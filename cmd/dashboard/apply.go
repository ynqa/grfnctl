package dashboard

import (
	"fmt"

	"github.com/grafana/grafana-openapi-client-go/models"
	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/internal/grafana"
	"github.com/ynqa/grfnctl/internal/util"
)

var (
	filePath  string
	folderUID string
	overwrite bool
)

func init() {
	ApplyCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the dashboard JSON file")
	ApplyCmd.RegisterFlagCompletionFunc(
		"file",
		cobra.FixedCompletions([]string{"json"}, cobra.ShellCompDirectiveFilterFileExt),
	)
	ApplyCmd.Flags().StringVarP(&folderUID, "folder", "F", "", "UID of the folder to save the dashboard in")
	ApplyCmd.RegisterFlagCompletionFunc(
		"folder",
		util.FolderCompletionFunc(),
	)
	ApplyCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing dashboard with the same UID")
	ApplyCmd.MarkFlagRequired("file")
}

var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Grafana dashboard from a JSON file",
	Long: `Apply Grafana dashboard to the configured Grafana instance from a JSON file.
The JSON file should contain the dashboard definition as per Grafana's API requirements.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := grafana.NewClientFromEnv()
		if err != nil {
			return err
		}

		dashboardJSON, err := util.ReadDashboardFile(filePath)
		if err != nil {
			return err
		}

		// NOTE: It has been confirmed that the dashboard update succeeds in the following cases:
		// - UID described in the JSON is the same as the existing dashboard
		// - `overwrite` flag is set to FALSE
		// Therefore, we check for UID conflicts on the grfnctl side in advance.
		dashboardUID, ok := dashboardJSON["uid"].(string)
		if ok {
			existed, err := grafana.GetDashboardBy(client, dashboardUID)
			if err == nil && existed.Payload != nil && !overwrite {
				return fmt.Errorf(
					"dashboard with UID '%s (%s)' already exists; use --overwrite to update",
					existed.Payload.Meta.Slug,
					dashboardUID,
				)
			}
		}

		params := &models.SaveDashboardCommand{
			Dashboard: dashboardJSON,
			FolderUID: folderUID,
			Overwrite: overwrite,
		}

		_, err = grafana.CreateOrUpdateDashboard(client, params)
		if err != nil {
			return err
		}

		return nil
	},
}
