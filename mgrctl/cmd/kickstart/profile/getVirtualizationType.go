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

type getVirtualizationTypeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getVirtualizationTypeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getVirtualizationType",
		Short: "For given kickstart profile label returns label of
 virtualization type it's using",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getVirtualizationTypeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getVirtualizationType)
		},
	}

	cmd.Flags().String("KsLabel", "", "")

	return cmd
}

func getVirtualizationType(globalFlags *types.GlobalFlags, flags *getVirtualizationTypeFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

