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

type setProfileNameFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
	Name                  string
}

func setProfileNameCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setProfileName",
		Short: "Set the profile name for the server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setProfileNameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setProfileName)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Name", "", "Name of the profile.")

	return cmd
}

func setProfileName(globalFlags *types.GlobalFlags, flags *setProfileNameFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
