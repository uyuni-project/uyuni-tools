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

type createProfileWithCustomUrlFlags struct {
	api.ConnectionDetails  `mapstructure:"api"`
	ProfileLabel           string
	VirtualizationType     string
	KickstartableTreeLabel string
	DownloadUrl            bool
	RootPassword           string
	UpdateType             string
}

func createProfileWithCustomUrlCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createProfileWithCustomUrl",
		Short: "Create a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createProfileWithCustomUrlFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createProfileWithCustomUrl)
		},
	}

	cmd.Flags().String("ProfileLabel", "", "Label for the new kickstart profile.")
	cmd.Flags().String("VirtualizationType", "", "none, para_host, qemu, xenfv or xenpv.")
	cmd.Flags().String("KickstartableTreeLabel", "", "Label of a kickstartable tree to associate the new profile with.")
	cmd.Flags().String("DownloadUrl", "", "Download URL, or 'default' to use the kickstart tree's default URL.")
	cmd.Flags().String("RootPassword", "", "Root password.")
	cmd.Flags().String("UpdateType", "", "Should the profile update itself to use the newest tree available? Possible values are: none (default) or all (includes custom Kickstart Trees).")

	return cmd
}

func createProfileWithCustomUrl(globalFlags *types.GlobalFlags, flags *createProfileWithCustomUrlFlags, cmd *cobra.Command, args []string) error {

	res, err := kickstart.Kickstart(&flags.ConnectionDetails, flags.ProfileLabel, flags.VirtualizationType, flags.KickstartableTreeLabel, flags.DownloadUrl, flags.RootPassword, flags.UpdateType)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
