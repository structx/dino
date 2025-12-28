package credentials

import (
	"github.com/spf13/cobra"
	"soft.structx.io/dino/cmd/cli/sub/tunnel"
)

func init() {
	credCmd.AddCommand(rotateCmd)

	tunnel.TunnelCmd.AddCommand(credCmd)
}

var (
	credCmd = &cobra.Command{
		Use:     "credentials",
		Aliases: []string{"c", "creds", "cred"},
	}
)
