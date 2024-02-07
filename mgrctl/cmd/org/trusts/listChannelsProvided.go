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

type listChannelsProvidedFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func listChannelsProvidedCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listChannelsProvided",
		Short: "Lists all software channels that the organization given is providing to
 the user's organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listChannelsProvidedFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listChannelsProvided)
		},
	}

	cmd.Flags().String("OrgId", "", "Id of the trusted organization")

	return cmd
}

func listChannelsProvided(globalFlags *types.GlobalFlags, flags *listChannelsProvidedFlags, cmd *cobra.Command, args []string) error {

res, err := trusts.Trusts(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

