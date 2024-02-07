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

type createOrUpdateSymlinkFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	Path          string
}

func createOrUpdateSymlinkCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createOrUpdateSymlink",
		Short: "Create a new symbolic link with the given path, or
 update an existing path in config channel of 'normal' type.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createOrUpdateSymlinkFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createOrUpdateSymlink)
		},
	}

	cmd.Flags().String("Label", "", "")
	cmd.Flags().String("Path", "", "")

	return cmd
}

func createOrUpdateSymlink(globalFlags *types.GlobalFlags, flags *createOrUpdateSymlinkFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.Path)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

