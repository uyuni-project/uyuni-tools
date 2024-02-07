package trusts

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/org/trusts"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listOrgsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listOrgsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listOrgs",
		Short: "List all organanizations trusted by the user's organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listOrgsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listOrgs)
		},
	}

	return cmd
}

func listOrgs(globalFlags *types.GlobalFlags, flags *listOrgsFlags, cmd *cobra.Command, args []string) error {

	res, err := trusts.Trusts(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
