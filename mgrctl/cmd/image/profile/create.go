package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/image/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	Type          string
	StoreLabel          string
	Path          string
	ActivationKey          string
	KiwiOptions          string
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new image profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("Label", "", "")
	cmd.Flags().String("Type", "", "")
	cmd.Flags().String("StoreLabel", "", "")
	cmd.Flags().String("Path", "", "")
	cmd.Flags().String("ActivationKey", "", "optional")
	cmd.Flags().String("KiwiOptions", "", "")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.Label, flags.Type, flags.StoreLabel, flags.Path, flags.ActivationKey, flags.KiwiOptions)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

