package route

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"soft.structx.io/dino/client"
	"soft.structx.io/dino/logging"
)

var (
	listCmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			logger := logging.FromContext(ctx)
			cli := client.FromContext(ctx)

			timeout, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			partials, err := cli.ListRoutes(timeout, "hello", 10, 0)
			if err != nil {
				return fmt.Errorf("cli.ListRoutes: %w", err)
			}

			for _, p := range partials {
				logger.Info("routes", zap.Any("partial", p))
			}

			return nil
		},
	}
)
