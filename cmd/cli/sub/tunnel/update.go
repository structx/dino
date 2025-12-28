package tunnel

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
	updateCmd = &cobra.Command{
		Use:   "update [NAME] [NEW_NAME]",
		Short: "update tunnel",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			logger := logging.FromContext(ctx)
			cli := client.FromContext(ctx)

			oldName := args[0]
			newName := args[1]
			if len(oldName) < 1 || len(newName) < 1 {
				return fmt.Errorf("missing name arg")
			}

			logger.Debug("update tunnel", zap.String("old_name", oldName), zap.String("new_name", newName))

			timeout, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			resp, err := cli.UpdateTunnel(timeout, client.TunnelUpdate{
				OldName: oldName,
				Name:    newName,
			})
			if err != nil {
				logger.Error("cli.UpdateTunnel", zap.Error(err))
				return fmt.Errorf("cli.UpdateTunnel: %w", err)
			}

			logger.Debug("update success", zap.Any("tunnel", resp))

			return nil
		},
	}
)
