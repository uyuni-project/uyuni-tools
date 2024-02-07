package custominfo

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/custominfo"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deleteKeyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KeyLabel              string
}

func deleteKeyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteKey",
		Short: "Delete an existing custom key and all systems' values for the key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteKeyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteKey)
		},
	}

	cmd.Flags().String("KeyLabel", "", "new key's label")

	return cmd
}

func deleteKey(globalFlags *types.GlobalFlags, flags *deleteKeyFlags, cmd *cobra.Command, args []string) error {

	res, err := custominfo.Custominfo(&flags.ConnectionDetails, flags.KeyLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
