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

type removeIpRangeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
	IpAddress             string
}

func removeIpRangeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeIpRange",
		Short: "Remove an ip range from a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeIpRangeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeIpRange)
		},
	}

	cmd.Flags().String("KsLabel", "", "The kickstart label of the ip range you want to remove")
	cmd.Flags().String("IpAddress", "", "An Ip Address that falls within the range that you are wanting to remove. The min or max of the range will work.")

	return cmd
}

func removeIpRange(globalFlags *types.GlobalFlags, flags *removeIpRangeFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.IpAddress)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
