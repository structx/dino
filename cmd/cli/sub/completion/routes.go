package completion

import "github.com/spf13/cobra"

// RouteHostnameFunc
func RouteHostnameFunc(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveDefault
}
