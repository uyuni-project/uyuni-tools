package kickstart

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAllIpRangesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAllIpRangesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllIpRanges",
		Short: "List all Ip Ranges and their associated kickstarts available
 in the user's org.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllIpRangesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllIpRanges)
		},
	}


	return cmd
}

func listAllIpRanges(globalFlags *types.GlobalFlags, flags *listAllIpRangesFlags, cmd *cobra.Command, args []string) error {

res, err := kickstart.Kickstart(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

