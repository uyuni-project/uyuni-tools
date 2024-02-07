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

type getCfgPreservationFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getCfgPreservationCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getCfgPreservation",
		Short: "Get ks.cfg preservation option for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getCfgPreservationFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getCfgPreservation)
		},
	}

	cmd.Flags().String("KsLabel", "", "Label of kickstart profile to be changed.")

	return cmd
}

func getCfgPreservation(globalFlags *types.GlobalFlags, flags *getCfgPreservationFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

