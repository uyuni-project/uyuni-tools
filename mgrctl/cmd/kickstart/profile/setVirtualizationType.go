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

type setVirtualizationTypeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	TypeLabel          string
}

func setVirtualizationTypeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setVirtualizationType",
		Short: "For given kickstart profile label sets its virtualization type.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setVirtualizationTypeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setVirtualizationType)
		},
	}

	cmd.Flags().String("KsLabel", "", "")
	cmd.Flags().String("TypeLabel", "", "One of the following: 'none', 'qemu', 'para_host', 'xenpv', 'xenfv'")

	return cmd
}

func setVirtualizationType(globalFlags *types.GlobalFlags, flags *setVirtualizationTypeFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.TypeLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

