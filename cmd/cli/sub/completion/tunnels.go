package completion

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"soft.structx.io/dino/client"
	"soft.structx.io/dino/logging"
)

// TunnelNameFunc
func TunnelNameFunc(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	ctx := cmd.Context()
	cli := client.FromContext(ctx)

	timeout, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	params := client.TunnelList{
		Limit:        defaultLimit,
		Offset:       defaultOffset,
		Autocomplete: true,
		CompleteMe:   toComplete,
	}

	results, err := cli.ListTunnels(timeout, params)
	if err != nil {
		logging.FromContext(ctx).Error("cli.ListTunnels", zap.Error(err))
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	finalCompletions := make([]cobra.Completion, 0, len(results))
	for _, r := range results {
		finalCompletions = append(finalCompletions, r.Name)
	}

	return finalCompletions, cobra.ShellCompDirectiveDefault
}
