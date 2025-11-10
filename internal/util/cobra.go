package util

import (
	"fmt"
	"strings"

	"github.com/grafana/grafana-openapi-client-go/client/snapshots"
	"github.com/spf13/cobra"

	"github.com/ynqa/grfnctl/internal/grafana"
)

// CompletionGuard allows callers to short-circuit completion when conditions are not met.
type CompletionGuard func(cmd *cobra.Command, args []string) (cobra.ShellCompDirective, bool)

// SkipCompletionWhenArgsProvided prevents suggesting dashboards when arguments already exist.
func SkipCompletionWhenArgsProvided() CompletionGuard {
	return func(cmd *cobra.Command, args []string) (cobra.ShellCompDirective, bool) {
		if len(args) > 0 {
			return cobra.ShellCompDirectiveNoFileComp, true
		}
		return cobra.ShellCompDirectiveDefault, false
	}
}

// DashboardCompletionFunc provides a completion function for dashboard UIDs.
func DashboardCompletionFunc(guards ...CompletionGuard) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		for _, guard := range guards {
			if directive, skip := guard(cmd, args); skip {
				return nil, directive
			}
		}

		client, err := grafana.NewClientFromEnv()
		if err != nil {
			cobra.CompError(err.Error())
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		searchResult, err := grafana.ListDashboards(client, nil)
		if err != nil {
			cobra.CompError(err.Error())
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		toComplete = strings.ToLower(toComplete)

		boards := grafana.FilterDashboardsByTypeDB(searchResult)

		suggestions := make([]string, 0, len(boards))
		for _, board := range boards {
			switch {
			case toComplete == "":
			case strings.HasPrefix(strings.ToLower(board.UID), toComplete):
			case strings.HasPrefix(strings.ToLower(board.Title), toComplete):
			default:
				continue
			}

			suggestions = append(suggestions, fmt.Sprintf("%s\t%s", board.UID, board.Title))
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}
}

// FolderCompletionFunc provides a completion function for folder UIDs.
func FolderCompletionFunc() cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := grafana.NewClientFromEnv()
		if err != nil {
			cobra.CompError(err.Error())
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		folders, err := grafana.ListFolders(client, nil)
		if err != nil {
			cobra.CompError(err.Error())
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		toComplete = strings.ToLower(toComplete)

		suggestions := make([]string, 0, len(folders.Payload))
		for _, folder := range folders.Payload {
			switch {
			case toComplete == "":
			case strings.HasPrefix(strings.ToLower(folder.UID), toComplete):
			case strings.HasPrefix(strings.ToLower(folder.Title), toComplete):
			default:
				continue
			}

			if folder.Title != "" {
				suggestions = append(suggestions, fmt.Sprintf("%s\t%s", folder.UID, folder.Title))
				continue
			}
			suggestions = append(suggestions, folder.UID)
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}
}

// SnapshotCompletionFunc provides a completion function for snapshot keys.
func SnapshotCompletionFunc(guards ...CompletionGuard) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		for _, guard := range guards {
			if directive, skip := guard(cmd, args); skip {
				return nil, directive
			}
		}

		client, err := grafana.NewClientFromEnv()
		if err != nil {
			cobra.CompError(err.Error())
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		snapshots, err := grafana.ListSnapshots(client, &snapshots.SearchDashboardSnapshotsParams{})
		if err != nil {
			cobra.CompError(err.Error())
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		toComplete = strings.ToLower(toComplete)
		suggestions := make([]string, 0, len(snapshots.Payload))
		for _, snapshot := range snapshots.Payload {
			switch {
			case toComplete == "":
			case strings.HasPrefix(strings.ToLower(snapshot.Key), toComplete):
			case strings.HasPrefix(strings.ToLower(snapshot.Name), toComplete):
			default:
				continue
			}

			if snapshot.Name != "" {
				suggestions = append(suggestions, fmt.Sprintf("%s\t%s", snapshot.Key, snapshot.Name))
				continue
			}
			suggestions = append(suggestions, snapshot.Key)
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}
}
