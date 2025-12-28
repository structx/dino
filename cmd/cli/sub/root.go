package sub

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	"github.com/spf13/cobra"
	"soft.structx.io/dino/client"
	"soft.structx.io/dino/logging"
)

var (
	targetFlag     string
	targetFlagName string = "target"

	RootCmd = &cobra.Command{
		Use: "dino",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			target, err := cmd.Flags().GetString(targetFlagName)
			if err != nil {
				return fmt.Errorf("failed to get persistent flag: %w", err)
			}

			if cli, err := client.New(
				client.WithTarget(target),
			); err != nil {
				return fmt.Errorf("client.New: %w", err)
			} else {
				ctx = client.WithContext(ctx, cli)
			}

			// TODO
			// add flag to control verbosity
			logger := logging.NewConsoleLogger("DEBUG")

			ctx = logging.WithContext(ctx, logger)

			cmd.SetContext(ctx)

			return nil
		},
	}
)

func init() {
	RootCmd.AddCommand(connectCmd)
	RootCmd.AddCommand(closeCmd)

	RootCmd.PersistentFlags().StringVarP(&targetFlag, targetFlagName, "t", "api.dino.docker:8000", "target addr")
}

func Execute() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := RootCmd.ExecuteContext(ctx); err != nil {
		_, _ = fmt.Fprint(os.Stdout, err)
	}
}

// CmdExecute cobra command execute helper function
func CmdExecute(t *testing.T, cmd *cobra.Command, args []string) (string, error) {
	t.Helper()
	ctx := t.Context()

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	cmd.SetContext(ctx)

	err := cmd.ExecuteContext(ctx)
	return strings.TrimSpace(buf.String()), err
}
