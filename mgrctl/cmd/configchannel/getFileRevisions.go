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

type getFileRevisionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
	FilePath              string
}

func getFileRevisionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getFileRevisions",
		Short: "Get list of revisions for specified config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getFileRevisionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getFileRevisions)
		},
	}

	cmd.Flags().String("Label", "", "label of config channel to lookup on")
	cmd.Flags().String("FilePath", "", "config file path to examine")

	return cmd
}

func getFileRevisions(globalFlags *types.GlobalFlags, flags *getFileRevisionsFlags, cmd *cobra.Command, args []string) error {

	res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.FilePath)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
