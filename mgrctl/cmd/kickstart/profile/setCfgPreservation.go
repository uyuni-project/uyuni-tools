package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setCfgPreservationFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
	Preserve              bool
}

func setCfgPreservationCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setCfgPreservation",
		Short: "Set ks.cfg preservation option for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setCfgPreservationFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setCfgPreservation)
		},
	}

	cmd.Flags().String("KsLabel", "", "Label of kickstart profile to be changed.")
	cmd.Flags().String("Preserve", "", "whether or not      ks.cfg and all %include fragments will be copied to /root.")

	return cmd
}

func setCfgPreservation(globalFlags *types.GlobalFlags, flags *setCfgPreservationFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.Preserve)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
