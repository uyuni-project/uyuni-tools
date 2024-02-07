package api

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/api"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getApiNamespacesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func getApiNamespacesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getApiNamespaces",
		Short: "Lists available API namespaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getApiNamespacesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getApiNamespaces)
		},
	}


	return cmd
}

func getApiNamespaces(globalFlags *types.GlobalFlags, flags *getApiNamespacesFlags, cmd *cobra.Command, args []string) error {

res, err := api.Api(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

