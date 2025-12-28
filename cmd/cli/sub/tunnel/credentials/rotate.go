package credentials

import "github.com/spf13/cobra"

var (
	rotateCmd = &cobra.Command{
		Use:     "rotate",
		Aliases: []string{"r"},
	}
)
