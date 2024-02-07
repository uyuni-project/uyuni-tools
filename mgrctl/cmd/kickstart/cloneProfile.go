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

type cloneProfileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabelToClone        string
	NewKsLabel            string
}

func cloneProfileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloneProfile",
		Short: "Clone a Kickstart Profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags cloneProfileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, cloneProfile)
		},
	}

	cmd.Flags().String("KsLabelToClone", "", "Label of the kickstart profile to clone")
	cmd.Flags().String("NewKsLabel", "", "label of the cloned profile")

	return cmd
}

func cloneProfile(globalFlags *types.GlobalFlags, flags *cloneProfileFlags, cmd *cobra.Command, args []string) error {

	res, err := kickstart.Kickstart(&flags.ConnectionDetails, flags.KsLabelToClone, flags.NewKsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
