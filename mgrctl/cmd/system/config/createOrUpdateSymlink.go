package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/config"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createOrUpdateSymlinkFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Path          string
}

func createOrUpdateSymlinkCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createOrUpdateSymlink",
		Short: "Create a new symbolic link with the given path, or
 update an existing path.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createOrUpdateSymlinkFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createOrUpdateSymlink)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Path", "", "the configuration file/directory path")

	return cmd
}

func createOrUpdateSymlink(globalFlags *types.GlobalFlags, flags *createOrUpdateSymlinkFlags, cmd *cobra.Command, args []string) error {

res, err := config.Config(&flags.ConnectionDetails, flags.Sid, flags.Path)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

