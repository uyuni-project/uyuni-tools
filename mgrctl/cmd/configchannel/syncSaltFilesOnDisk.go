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

type syncSaltFilesOnDiskFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	$param.getFlagName()          $param.getType()
}

func syncSaltFilesOnDiskCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "syncSaltFilesOnDisk",
		Short: "Synchronize all files on the disk to the current state of the database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags syncSaltFilesOnDiskFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, syncSaltFilesOnDisk)
		},
	}

	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func syncSaltFilesOnDisk(globalFlags *types.GlobalFlags, flags *syncSaltFilesOnDiskFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

