package tunnel

import (
	"github.com/spf13/cobra"
	"soft.structx.io/dino/cmd/cli/sub"
)

var (
	nameFlag string

	limitFlag  uint32
	offsetFlag uint32

	tunnelIDFlag     string
	tunnelIDFlagName string = "tunnel-id"

	TunnelCmd = &cobra.Command{
		Use:   "tunnel",
		Short: "tunnel command group",
	}
)

func init() {
	addCmd.Flags().StringVar(&nameFlag, "name", "", "tunnel name")

	getCmd.Flags().StringVar(&tunnelIDFlag, tunnelIDFlagName, "", "tunnel id")

	listCmd.Flags().Uint32Var(&limitFlag, "limit", 10, "limit response items")
	listCmd.Flags().Uint32Var(&offsetFlag, "offset", 0, "offset response items")

	delCmd.Flags().StringVar(&tunnelIDFlag, tunnelIDFlagName, "", "tunnel id")

	TunnelCmd.AddCommand(addCmd)
	TunnelCmd.AddCommand(getCmd)
	TunnelCmd.AddCommand(listCmd)
	TunnelCmd.AddCommand(updateCmd)
	TunnelCmd.AddCommand(delCmd)

	sub.RootCmd.AddCommand(TunnelCmd)
}
