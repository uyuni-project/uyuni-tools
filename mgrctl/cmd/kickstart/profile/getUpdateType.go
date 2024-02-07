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

type getUpdateTypeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getUpdateTypeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getUpdateType",
		Short: "Get the update type for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getUpdateTypeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getUpdateType)
		},
	}

	cmd.Flags().String("KsLabel", "", "Label of kickstart profile.")

	return cmd
}

func getUpdateType(globalFlags *types.GlobalFlags, flags *getUpdateTypeFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

