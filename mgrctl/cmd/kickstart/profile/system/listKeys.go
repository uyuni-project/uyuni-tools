package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listKeysFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func listKeysCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listKeys",
		Short: "Returns the set of all keys associated with the given kickstart
             profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listKeysFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listKeys)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func listKeys(globalFlags *types.GlobalFlags, flags *listKeysFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

