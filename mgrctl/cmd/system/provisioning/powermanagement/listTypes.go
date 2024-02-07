package powermanagement

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/provisioning/powermanagement"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listTypesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listTypesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listTypes",
		Short: "Return a list of available power management types",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listTypesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listTypes)
		},
	}


	return cmd
}

func listTypes(globalFlags *types.GlobalFlags, flags *listTypesFlags, cmd *cobra.Command, args []string) error {

res, err := powermanagement.Powermanagement(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

