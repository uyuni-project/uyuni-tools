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

type getNetworkFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func getNetworkCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getNetwork",
		Short: "Get the addresses and hostname for a given server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getNetworkFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getNetwork)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getNetwork(globalFlags *types.GlobalFlags, flags *getNetworkFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

