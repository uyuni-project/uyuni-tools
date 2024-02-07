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

type appendToSoftwareListFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	$param.getFlagName()          $param.getType()
}

func appendToSoftwareListCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "appendToSoftwareList",
		Short: "Append the list of software packages to a kickstart profile.
 Duplicate packages will be ignored.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags appendToSoftwareListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, appendToSoftwareList)
		},
	}

	cmd.Flags().String("KsLabel", "", "the label of the kickstart profile")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func appendToSoftwareList(globalFlags *types.GlobalFlags, flags *appendToSoftwareListFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.KsLabel, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

