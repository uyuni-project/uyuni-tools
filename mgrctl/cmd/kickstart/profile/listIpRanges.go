package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listIpRangesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
}

func listIpRangesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listIpRanges",
		Short: "List all ip ranges for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listIpRangesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listIpRanges)
		},
	}

	cmd.Flags().String("KsLabel", "", "The label of the kickstart")

	return cmd
}

func listIpRanges(globalFlags *types.GlobalFlags, flags *listIpRangesFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
