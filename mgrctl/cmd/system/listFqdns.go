package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listFqdnsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func listFqdnsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listFqdns",
		Short: "Provides a list of FQDNs associated with a system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFqdnsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listFqdns)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listFqdns(globalFlags *types.GlobalFlags, flags *listFqdnsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
