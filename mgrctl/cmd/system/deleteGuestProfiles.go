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

type deleteGuestProfilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	HostId                int
	GuestNames            []string
}

func deleteGuestProfilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteGuestProfiles",
		Short: "Delete the specified list of guest profiles for a given host",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteGuestProfilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteGuestProfiles)
		},
	}

	cmd.Flags().String("HostId", "", "")
	cmd.Flags().String("GuestNames", "", "$desc")

	return cmd
}

func deleteGuestProfiles(globalFlags *types.GlobalFlags, flags *deleteGuestProfilesFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.HostId, flags.GuestNames)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
