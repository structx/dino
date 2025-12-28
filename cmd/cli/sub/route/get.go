package route

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"soft.structx.io/dino/client"
	"soft.structx.io/dino/cmd/cli/sub/completion"
	"soft.structx.io/dino/logging"
)

var (
	getCmd = &cobra.Command{
		Use:               "get [HOSTNAME]",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.RouteHostnameFunc,
		Short:             "get route",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			hostname := args[0]
			if len(hostname) < 1 {
				return fmt.Errorf("unexpected hostname length: %d", len(hostname))
			}

			logger := logging.FromContext(ctx)
			cli := client.FromContext(ctx)

			timeout, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			route, err := cli.GetRoute(timeout, hostname)
			if err != nil {
				return fmt.Errorf("cli.GetRoute: %w", err)
			}

			logger.Info("get route", zap.Any("details", route))

			return nil
		},
	}
)
