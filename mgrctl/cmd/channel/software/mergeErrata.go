package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type mergeErrataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MergeFromLabel          string
	MergeToLabel          string
	StartDate          string
	EndDate          string
	$param.getFlagName()          $param.getType()
}

func mergeErrataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mergeErrata",
		Short: "Merges all errata from one channel into another",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags mergeErrataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, mergeErrata)
		},
	}

	cmd.Flags().String("MergeFromLabel", "", "the label of the channel to pull errata from")
	cmd.Flags().String("MergeToLabel", "", "the label to push the errata into")
	cmd.Flags().String("StartDate", "", "")
	cmd.Flags().String("EndDate", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func mergeErrata(globalFlags *types.GlobalFlags, flags *mergeErrataFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.MergeFromLabel, flags.MergeToLabel, flags.StartDate, flags.EndDate, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

