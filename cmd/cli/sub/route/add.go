package route

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"soft.structx.io/dino/client"
	"soft.structx.io/dino/cmd/cli/sub/completion"
	"soft.structx.io/dino/logging"
)

const (
	hostnameFlagName string = "hostname"
	protocolFlagName string = "protocol"
	addrFlagName     string = "address"
	tunnelFlagName   string = "tunnel"
)

var (
	hostnameFlag string
	protocolFlag string
	addrFlag     string
	tunnelFlag   string
)

func init() {
	addCmd.Flags().StringVarP(&hostnameFlag, hostnameFlagName, "p", "", "public hostname (echo.dino.local)")
	addCmd.Flags().StringVarP(&protocolFlag, protocolFlagName, "r", "", "route local protocol (http)")
	addCmd.Flags().StringVarP(&addrFlag, addrFlagName, "a", "", "route local addr (localhost:8080)")
	addCmd.Flags().StringVarP(&tunnelFlag, tunnelFlagName, "x", "", "tunnel flag name (K3D_01)")

	_ = addCmd.MarkFlagRequired(hostnameFlagName)
	_ = addCmd.MarkFlagRequired(protocolFlagName)
	_ = addCmd.MarkFlagRequired(addrFlagName)
	_ = addCmd.MarkFlagRequired(tunnelFlagName)

	_ = addCmd.RegisterFlagCompletionFunc(tunnelFlagName, completion.TunnelNameFunc)
}

var (
	addCmd = &cobra.Command{
		Use: "add",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			logger := logging.FromContext(ctx)
			cli := client.FromContext(ctx)

			host, err := cmd.Flags().GetString(hostnameFlagName)
			if err != nil || len(host) < 1 {
				return fmt.Errorf("missing or unexpected host value %s, err: %w", host, err)
			}

			protocol, err := cmd.Flags().GetString(protocolFlagName)
			if err != nil || len(protocol) < 1 {
				return fmt.Errorf("missing or unexpected protocol value %s, err: %w", protocol, err)
			}

			addrFlag, err := cmd.Flags().GetString(addrFlagName)
			if err != nil || len(addrFlag) < 1 {
				return fmt.Errorf("missing or unexpected addr value %s, err: %w", addrFlag, err)
			}

			tunnel, err := cmd.Flags().GetString(tunnelFlagName)
			if err != nil || len(tunnel) < 1 {
				return fmt.Errorf("missing or unexpected tunnel value %s, err :%w", tunnel, err)
			}

			localHost, localPort, err := net.SplitHostPort(addrFlag)
			if err != nil {
				return fmt.Errorf("net.SplitHostPort: %w", err)
			}

			portU32, err := strconv.ParseUint(localPort, 10, 32)
			if err != nil {
				return fmt.Errorf("strconv.ParseUint: %w", err)
			}

			timeout, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			params := client.RouteAdd{
				Hostname:            host,
				Tunnel:              tunnel,
				DestinationProtocol: protocol,
				DestinationIP:       localHost,
				DestinationPort:     uint32(portU32),
				Enabled:             true,
			}

			logger.Debug("add route", zap.Any("params", params))

			route, err := cli.AddRoute(timeout, params)
			if err != nil {
				return fmt.Errorf("cli.AddRoute: %w", err)
			}

			logger.Info("success add route", zap.Any("new_route", route))

			return nil
		},
	}
)
