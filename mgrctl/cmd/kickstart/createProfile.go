package kickstart

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createProfileFlags struct {
	api.ConnectionDetails  `mapstructure:"api"`
	ProfileLabel           string
	VirtualizationType     string
	KickstartableTreeLabel string
	KickstartHost          string
	RootPassword           string
	UpdateType             string
}

func createProfileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createProfile",
		Short: "Create a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createProfileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createProfile)
		},
	}

	cmd.Flags().String("ProfileLabel", "", "Label for the new kickstart profile.")
	cmd.Flags().String("VirtualizationType", "", "none, para_host, qemu, xenfv or xenpv.")
	cmd.Flags().String("KickstartableTreeLabel", "", "Label of a kickstartable tree to associate the new profile with.")
	cmd.Flags().String("KickstartHost", "", "Kickstart hostname (of a #product() server or proxy) used to construct the default download URL for the new kickstart profile.")
	cmd.Flags().String("RootPassword", "", "Root password.")
	cmd.Flags().String("UpdateType", "", "Should the profile update itself to use the newest tree available? Possible values are: none (default) or all (includes custom Kickstart Trees).")

	return cmd
}

func createProfile(globalFlags *types.GlobalFlags, flags *createProfileFlags, cmd *cobra.Command, args []string) error {

	res, err := kickstart.Kickstart(&flags.ConnectionDetails, flags.ProfileLabel, flags.VirtualizationType, flags.KickstartableTreeLabel, flags.KickstartHost, flags.RootPassword, flags.UpdateType)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
