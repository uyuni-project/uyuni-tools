package errata

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/errata"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAffectedSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
}

func listAffectedSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAffectedSystems",
		Short: "Return the list of systems affected by the errata with the given advisory name.
 For those errata that are present in both vendor and user organizations under the same advisory name,
 this method retrieves the affected systems by both of them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAffectedSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAffectedSystems)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")

	return cmd
}

func listAffectedSystems(globalFlags *types.GlobalFlags, flags *listAffectedSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

