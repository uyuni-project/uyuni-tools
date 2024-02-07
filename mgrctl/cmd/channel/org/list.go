package org

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/org"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
}

func listCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the organizations associated with the given channel
 that may be trusted.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, list)
		},
	}

	cmd.Flags().String("Label", "", "label of the channel")

	return cmd
}

func list(globalFlags *types.GlobalFlags, flags *listFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

