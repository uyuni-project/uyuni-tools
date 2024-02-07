package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getSoftwareListFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getSoftwareListCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSoftwareList",
		Short: "Get a list of a kickstart profile's software packages.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSoftwareListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSoftwareList)
		},
	}

	cmd.Flags().String("KsLabel", "", "the label of the kickstart profile")

	return cmd
}

func getSoftwareList(globalFlags *types.GlobalFlags, flags *getSoftwareListFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

