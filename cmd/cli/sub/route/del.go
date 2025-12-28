package route

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
		Use:               "delete [HOSTNAME]",
		Aliases:           []string{"del"},
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.RouteHostnameFunc,
		Short:             "delete route",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := cmd.Context()

			hostname := args[0]
			if len(hostname) < 1 {
				return fmt.Errorf("unexpected host length: %d", len(hostname))
			}

			logger := logging.FromContext(ctx)
			cli := client.FromContext(ctx)

			timeout, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			err := cli.DelRoute(timeout, hostname)
			if err != nil {
				return fmt.Errorf("cli.DelRoute: %w", err)
			}

			logger.Info("success")

			return nil
		},
	}
)
