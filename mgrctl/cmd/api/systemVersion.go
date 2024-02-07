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

type systemVersionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func systemVersionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "systemVersion",
		Short: "Returns the server version.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags systemVersionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, systemVersion)
		},
	}

	return cmd
}

func systemVersion(globalFlags *types.GlobalFlags, flags *systemVersionFlags, cmd *cobra.Command, args []string) error {

	res, err := api.Api(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
