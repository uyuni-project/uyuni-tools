package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listGlobalsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listGlobalsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listGlobals",
		Short: "List all the global config channels accessible to the logged-in user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listGlobalsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listGlobals)
		},
	}

	return cmd
}

func listGlobals(globalFlags *types.GlobalFlags, flags *listGlobalsFlags, cmd *cobra.Command, args []string) error {

	res, err := configchannel.Configchannel(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
