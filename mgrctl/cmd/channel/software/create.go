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

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	Name          string
	Summary          string
	ArchLabel          string
	ParentLabel          string
	$param.getFlagName()          $param.getType()
	GpgCheck          bool
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a software channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("Label", "", "label of the new channel")
	cmd.Flags().String("Name", "", "name of the new channel")
	cmd.Flags().String("Summary", "", "summary of the channel")
	cmd.Flags().String("ArchLabel", "", "the label of the architecture the channel corresponds to,              run channel.software.listArches API for complete listing")
	cmd.Flags().String("ParentLabel", "", "label of the parent of this              channel, an empty string if it does not have one")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("GpgCheck", "", "true if the GPG check should be     enabled by default, false otherwise")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Label, flags.Name, flags.Summary, flags.ArchLabel, flags.ParentLabel, flags.$param.getFlagName(), flags.GpgCheck)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

