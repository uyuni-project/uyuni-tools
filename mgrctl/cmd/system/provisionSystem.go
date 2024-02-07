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

type provisionSystemFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	ProfileName          string
	EarliestDate          $date
}

func provisionSystemCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provisionSystem",
		Short: "Provision a system using the specified kickstart/autoinstallation profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags provisionSystemFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, provisionSystem)
		},
	}

	cmd.Flags().String("Sid", "", "ID of the system to be provisioned.")
	cmd.Flags().String("ProfileName", "", "Profile to use.")
	cmd.Flags().String("EarliestDate", "", "")

	return cmd
}

func provisionSystem(globalFlags *types.GlobalFlags, flags *provisionSystemFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.ProfileName, flags.EarliestDate)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

