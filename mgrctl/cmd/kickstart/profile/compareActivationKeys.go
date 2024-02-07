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

type compareActivationKeysFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KickstartLabel1          string
	KickstartLabel2          string
}

func compareActivationKeysCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compareActivationKeys",
		Short: "Returns a list for each kickstart profile; each list will contain
             activation keys not present on the other profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags compareActivationKeysFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, compareActivationKeys)
		},
	}

	cmd.Flags().String("KickstartLabel1", "", "")
	cmd.Flags().String("KickstartLabel2", "", "")

	return cmd
}

func compareActivationKeys(globalFlags *types.GlobalFlags, flags *compareActivationKeysFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KickstartLabel1, flags.KickstartLabel2)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

