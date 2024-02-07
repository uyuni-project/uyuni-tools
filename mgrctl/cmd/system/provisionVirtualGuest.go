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

type provisionVirtualGuestFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	GuestName          string
	ProfileName          string
	ProfileName          string
	MemoryMb          int
	Vcpus          int
	StorageGb          int
	MacAddress          string
}

func provisionVirtualGuestCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provisionVirtualGuest",
		Short: "Provision a guest on the host specified.  Defaults to:
 memory=512MB, vcpu=1, storage=3GB, mac_address=random.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags provisionVirtualGuestFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, provisionVirtualGuest)
		},
	}

	cmd.Flags().String("Sid", "", "ID of host to provision guest on.")
	cmd.Flags().String("GuestName", "", "")
	cmd.Flags().String("ProfileName", "", "Kickstart profile to use.")
	cmd.Flags().String("ProfileName", "", "Kickstart Profile to use.")
	cmd.Flags().String("MemoryMb", "", "Memory to allocate to the guest")
	cmd.Flags().String("Vcpus", "", "Number of virtual CPUs to allocate to                                          the guest.")
	cmd.Flags().String("StorageGb", "", "Size of the guests disk image.")
	cmd.Flags().String("MacAddress", "", "macAddress to give the guest's                                          virtual networking hardware.")

	return cmd
}

func provisionVirtualGuest(globalFlags *types.GlobalFlags, flags *provisionVirtualGuestFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.GuestName, flags.ProfileName, flags.ProfileName, flags.MemoryMb, flags.Vcpus, flags.StorageGb, flags.MacAddress)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

