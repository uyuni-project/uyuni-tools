package utils

import (
	"github.com/spf13/cobra"
)

func AddBackendFlag(cmd *cobra.Command) {
	cmd.Flags().String("backend", "", "tool to use to reach the container. Possible values: 'podman', 'kubectl'. Default guesses which to use.")
}
