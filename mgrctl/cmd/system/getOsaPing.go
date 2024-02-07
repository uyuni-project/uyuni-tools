package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getOsaPingFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	LoggedInUser          User
	Sid                   int
}

func getOsaPingCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getOsaPing",
		Short: "get details about a ping sent to a system using OSA",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getOsaPingFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getOsaPing)
		},
	}

	cmd.Flags().String("LoggedInUser", "", "")
	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getOsaPing(globalFlags *types.GlobalFlags, flags *getOsaPingFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.LoggedInUser, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
