package tunnel

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
		Use:               "get [NAME]",
		Short:             "get tunnel",
		ValidArgsFunction: completion.TunnelNameFunc,
		Args:              cobra.ExactArgs(1),
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

			tunnel, err := cli.GetTunnel(timeout, name)
			if err != nil {
				logger.Error("cli.GetTunnel", zap.Error(err))
				return fmt.Errorf("cli.GetTunnel: %w", err)
			}

			logger.Info("tunnel", zap.Any("tunnel", tunnel))

			return nil
		},
	}
)
