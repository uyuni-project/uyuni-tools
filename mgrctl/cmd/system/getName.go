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

type getNameFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          string
}

func getNameCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getName",
		Short: "Get system name and last check in information for the given system ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getNameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getName)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getName(globalFlags *types.GlobalFlags, flags *getNameFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

