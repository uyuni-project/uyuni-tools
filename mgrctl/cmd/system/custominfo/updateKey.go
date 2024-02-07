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

type updateKeyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KeyLabel          string
	KeyDescription          string
}

func updateKeyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateKey",
		Short: "Update description of a custom key",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateKeyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateKey)
		},
	}

	cmd.Flags().String("KeyLabel", "", "key to change")
	cmd.Flags().String("KeyDescription", "", "new key's description")

	return cmd
}

func updateKey(globalFlags *types.GlobalFlags, flags *updateKeyFlags, cmd *cobra.Command, args []string) error {

res, err := custominfo.Custominfo(&flags.ConnectionDetails, flags.KeyLabel, flags.KeyDescription)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

