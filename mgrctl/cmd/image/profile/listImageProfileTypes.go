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

type listImageProfileTypesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listImageProfileTypesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listImageProfileTypes",
		Short: "List available image store types",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listImageProfileTypesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listImageProfileTypes)
		},
	}

	return cmd
}

func listImageProfileTypes(globalFlags *types.GlobalFlags, flags *listImageProfileTypesFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
