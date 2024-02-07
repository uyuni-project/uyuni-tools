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

type compareAdvancedOptionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KickstartLabel1          string
	KickstartLabel2          string
}

func compareAdvancedOptionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compareAdvancedOptions",
		Short: "Returns a list for each kickstart profile; each list will contain the
             properties that differ between the profiles and their values for that
             specific profile .",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags compareAdvancedOptionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, compareAdvancedOptions)
		},
	}

	cmd.Flags().String("KickstartLabel1", "", "")
	cmd.Flags().String("KickstartLabel2", "", "")

	return cmd
}

func compareAdvancedOptions(globalFlags *types.GlobalFlags, flags *compareAdvancedOptionsFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KickstartLabel1, flags.KickstartLabel2)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

