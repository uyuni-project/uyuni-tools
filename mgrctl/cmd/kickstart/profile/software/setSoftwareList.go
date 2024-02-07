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

type setSoftwareListFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	$param.getFlagName()          $param.getType()
	$param.getFlagName()          $param.getType()
	IgnoreMissing          bool
	NoBase          bool
}

func setSoftwareListCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setSoftwareList",
		Short: "Set the list of software packages for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setSoftwareListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setSoftwareList)
		},
	}

	cmd.Flags().String("KsLabel", "", "the label of the kickstart profile")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("IgnoreMissing", "", "ignore missing packages if true")
	cmd.Flags().String("NoBase", "", "don't install @Base package group if true")

	return cmd
}

func setSoftwareList(globalFlags *types.GlobalFlags, flags *setSoftwareListFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.KsLabel, flags.$param.getFlagName(), flags.$param.getFlagName(), flags.IgnoreMissing, flags.NoBase)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

