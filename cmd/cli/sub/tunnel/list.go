package tunnel

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"soft.structx.io/dino/client"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list tunnels",
		// ValidArgs: []cobra.Completion{limitFlag, offsetFlag},
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := cmd.Context()
			cli := client.FromContext(ctx)

			timeout, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			partials, err := cli.ListTunnels(timeout, client.TunnelList{
				Limit:        10,
				Offset:       0,
				Autocomplete: false,
				CompleteMe:   "",
			})
			if err != nil {
				return fmt.Errorf("cli.ListTunnels: %w", err)
			}

			fmt.Printf("found %d tunnels\n", len(partials))
			for _, p := range partials {
				fmt.Println("----")
				fmt.Println(p)
			}

			return nil
		},
	}
)
