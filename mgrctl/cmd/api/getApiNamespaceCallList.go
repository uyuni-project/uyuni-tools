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

type getApiNamespaceCallListFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Namespace             string
}

func getApiNamespaceCallListCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getApiNamespaceCallList",
		Short: "Lists all available api calls for the specified namespace",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getApiNamespaceCallListFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getApiNamespaceCallList)
		},
	}

	cmd.Flags().String("Namespace", "", "")

	return cmd
}

func getApiNamespaceCallList(globalFlags *types.GlobalFlags, flags *getApiNamespaceCallListFlags, cmd *cobra.Command, args []string) error {

	res, err := api.Api(&flags.ConnectionDetails, flags.Namespace)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
