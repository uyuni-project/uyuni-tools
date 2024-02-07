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

type deleteFilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
	Paths                 []string
}

func deleteFilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteFiles",
		Short: "Remove file paths from a global channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteFilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteFiles)
		},
	}

	cmd.Flags().String("Label", "", "channel to remove the files from")
	cmd.Flags().String("Paths", "", "$desc")

	return cmd
}

func deleteFiles(globalFlags *types.GlobalFlags, flags *deleteFilesFlags, cmd *cobra.Command, args []string) error {

	res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.Paths)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
