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

type createOrUpdatePathFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Path          string
}

func createOrUpdatePathCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createOrUpdatePath",
		Short: "Create a new file (text or binary) or directory with the given path, or
 update an existing path on a server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createOrUpdatePathFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createOrUpdatePath)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Path", "", "the configuration file/directory path")

	return cmd
}

func createOrUpdatePath(globalFlags *types.GlobalFlags, flags *createOrUpdatePathFlags, cmd *cobra.Command, args []string) error {

res, err := config.Config(&flags.ConnectionDetails, flags.Sid, flags.Path)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

