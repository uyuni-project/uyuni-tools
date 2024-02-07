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

type lookupFileInfoFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	$param.getFlagName()          $param.getType()
}

func lookupFileInfoCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookupFileInfo",
		Short: "Given a list of paths and a server, returns details about
 the latest revisions of the paths.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags lookupFileInfoFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, lookupFileInfo)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func lookupFileInfo(globalFlags *types.GlobalFlags, flags *lookupFileInfoFlags, cmd *cobra.Command, args []string) error {

res, err := config.Config(&flags.ConnectionDetails, flags.Sid, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

