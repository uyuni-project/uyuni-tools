package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type lookupFileInfoFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	$param.getFlagName()          $param.getType()
	Path          string
	Revision          int
}

func lookupFileInfoCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookupFileInfo",
		Short: "Given a list of paths and a channel, returns details about
 the latest revisions of the paths.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags lookupFileInfoFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, lookupFileInfo)
		},
	}

	cmd.Flags().String("Label", "", "label of config channel to lookup on")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("Path", "", "path of file/directory")
	cmd.Flags().String("Revision", "", "the revision number")

	return cmd
}

func lookupFileInfo(globalFlags *types.GlobalFlags, flags *lookupFileInfoFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.$param.getFlagName(), flags.Path, flags.Revision)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

