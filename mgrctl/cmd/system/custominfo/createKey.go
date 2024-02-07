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

type createKeyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KeyLabel              string
	KeyDescription        string
}

func createKeyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createKey",
		Short: "Create a new custom key",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createKeyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createKey)
		},
	}

	cmd.Flags().String("KeyLabel", "", "new key's label")
	cmd.Flags().String("KeyDescription", "", "new key's description")

	return cmd
}

func createKey(globalFlags *types.GlobalFlags, flags *createKeyFlags, cmd *cobra.Command, args []string) error {

	res, err := custominfo.Custominfo(&flags.ConnectionDetails, flags.KeyLabel, flags.KeyDescription)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
