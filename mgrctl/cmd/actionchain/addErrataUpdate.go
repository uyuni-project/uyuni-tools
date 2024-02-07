package actionchain

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/actionchain"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addErrataUpdateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	$param.getFlagName()          $param.getType()
	ChainLabel          string
	$param.getFlagName()          $param.getType()
	OnlyRelevant          bool
}

func addErrataUpdateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addErrataUpdate",
		Short: "Adds Errata update to an Action Chain.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addErrataUpdateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addErrataUpdate)
		},
	}

	cmd.Flags().String("Sid", "", "System ID")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("ChainLabel", "", "Label of the chain")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("OnlyRelevant", "", "If true, InvalidErrataException is thrown if an errata does not apply to a system.")

	return cmd
}

func addErrataUpdate(globalFlags *types.GlobalFlags, flags *addErrataUpdateFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.Sid, flags.$param.getFlagName(), flags.ChainLabel, flags.$param.getFlagName(), flags.OnlyRelevant)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

