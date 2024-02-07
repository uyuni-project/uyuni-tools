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

type listFilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
}

func listFilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listFiles",
		Short: "Return a list of files in a channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listFiles)
		},
	}

	cmd.Flags().String("Label", "", "label of config channel to list files on")

	return cmd
}

func listFiles(globalFlags *types.GlobalFlags, flags *listFilesFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

