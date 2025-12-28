package route

import (
	"github.com/spf13/cobra"
	"soft.structx.io/dino/cmd/cli/sub"
)

var (
	routeCmd = &cobra.Command{
		Use:   "route",
		Short: "route command group",
	}
)

func init() {
	routeCmd.AddCommand(addCmd)
	routeCmd.AddCommand(delCmd)
	routeCmd.AddCommand(getCmd)
	routeCmd.AddCommand(listCmd)
	routeCmd.AddCommand(updateCmd)

	sub.RootCmd.AddCommand(routeCmd)
}
