package tunnel

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"soft.structx.io/dino/client"
	"soft.structx.io/dino/logging"
)

var (
	addCmd = &cobra.Command{
		Use:   "add [NAME]",
		Short: "add tunnel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			name := args[0]
			if len(name) < 1 {
				return fmt.Errorf("unexpected name length %d", len(name))
			}

			logger := logging.FromContext(ctx)
			cli := client.FromContext(ctx)

			// timeout, cancel := context.WithTimeout(ctx, time.Second*15)
			// defer cancel()

			tunnelName := strings.TrimSpace(name)
			tunnelName = strings.ToValidUTF8(tunnelName, "")

			logger.Debug("add tunnel", zap.String("tunnel_name", tunnelName))

			resp, auth, err := cli.AddTunnel(ctx, client.TunnelAdd{Name: tunnelName})
			if err != nil {
				return fmt.Errorf("cli.AddTunnel: %w", err)
			}

			logger.Debug("tunnel add successful", zap.Any("tunnel", resp))
			logger.Debug("tunnel credentials", zap.Any("auth_details", auth))

			return nil
		},
	}
)
