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

type importFileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProfileLabel          string
	VirtualizationType          string
	KickstartableTreeLabel          string
	KickstartFileContents          string
	KickstartHost          string
	UpdateType          string
}

func importFileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "importFile",
		Short: "Import a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags importFileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, importFile)
		},
	}

	cmd.Flags().String("ProfileLabel", "", "Label for the new kickstart profile.")
	cmd.Flags().String("VirtualizationType", "", "none, para_host, qemu, xenfv or xenpv.")
	cmd.Flags().String("KickstartableTreeLabel", "", "Label of a kickstartable tree to associate the new profile with.")
	cmd.Flags().String("KickstartFileContents", "", "Contents of the kickstart file to import.")
	cmd.Flags().String("KickstartHost", "", "Kickstart hostname (of a #product() server or proxy) used to construct the default download URL for the new kickstart profile. Using this option signifies that this default URL will be used instead of any url/nfs/cdrom/harddrive commands in the kickstart file itself.")
	cmd.Flags().String("UpdateType", "", "Should the profile update itself to use the newest tree available? Possible values are: none (default) or all (includes custom Kickstart Trees).")

	return cmd
}

func importFile(globalFlags *types.GlobalFlags, flags *importFileFlags, cmd *cobra.Command, args []string) error {

res, err := kickstart.Kickstart(&flags.ConnectionDetails, flags.ProfileLabel, flags.VirtualizationType, flags.KickstartableTreeLabel, flags.KickstartFileContents, flags.KickstartHost, flags.UpdateType)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

