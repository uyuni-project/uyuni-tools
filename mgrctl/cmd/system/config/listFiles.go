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

type listFilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func listFilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listFiles",
		Short: "Return the list of files in a given channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listFiles)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listFiles(globalFlags *types.GlobalFlags, flags *listFilesFlags, cmd *cobra.Command, args []string) error {

	res, err := config.Config(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
