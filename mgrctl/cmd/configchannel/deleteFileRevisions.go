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

type deleteFileRevisionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	FilePath          string
	$param.getFlagName()          $param.getType()
}

func deleteFileRevisionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteFileRevisions",
		Short: "Delete specified revisions of a given configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteFileRevisionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteFileRevisions)
		},
	}

	cmd.Flags().String("Label", "", "label of config channel to lookup on")
	cmd.Flags().String("FilePath", "", "configuration file path")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func deleteFileRevisions(globalFlags *types.GlobalFlags, flags *deleteFileRevisionsFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.FilePath, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

