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

type deleteFilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
	Paths                 []string
}

func deleteFilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteFiles",
		Short: "Removes file paths from a local or sandbox channel of a server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteFilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteFiles)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Paths", "", "$desc")

	return cmd
}

func deleteFiles(globalFlags *types.GlobalFlags, flags *deleteFilesFlags, cmd *cobra.Command, args []string) error {

	res, err := config.Config(&flags.ConnectionDetails, flags.Sid, flags.Paths)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
