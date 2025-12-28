package tunnel

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"soft.structx.io/dino/client"
	"soft.structx.io/dino/cmd/cli/sub/completion"
	"soft.structx.io/dino/logging"
)

var (
	delCmd = &cobra.Command{
		Use:               "delete [NAME]",
		Aliases:           []string{"del"},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.TunnelNameFunc,
		Short:             "delete tunnel",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := cmd.Context()

			name := args[0]
			if len(name) < 1 {
				return fmt.Errorf("unexpected name length: %d", len(name))
			}

			logger := logging.FromContext(ctx)
			cli := client.FromContext(ctx)

			timeout, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			err := cli.DelTunnel(timeout, name)
			if err != nil {
				return fmt.Errorf("cli.DelTunnel: %w", err)
			}

			logger.Info("success")

			return nil
		},
	}
)
