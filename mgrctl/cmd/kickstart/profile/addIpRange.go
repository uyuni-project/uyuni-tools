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

type addIpRangeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
	Min                   string
	Max                   string
}

func addIpRangeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addIpRange",
		Short: "Add an ip range to a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addIpRangeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addIpRange)
		},
	}

	cmd.Flags().String("KsLabel", "", "The label of the kickstart")
	cmd.Flags().String("Min", "", "The ip address making up the minimum of the range (i.e. 192.168.0.1)")
	cmd.Flags().String("Max", "", "The ip address making up the maximum of the range (i.e. 192.168.0.254)")

	return cmd
}

func addIpRange(globalFlags *types.GlobalFlags, flags *addIpRangeFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.Min, flags.Max)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
