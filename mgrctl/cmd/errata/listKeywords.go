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

type listKeywordsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
}

func listKeywordsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listKeywords",
		Short: "Get the keywords associated with an erratum matching the given advisory name.
 For those errata that are present in both vendor and user organizations under the same advisory name,
 this method retrieves the keywords of both of them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listKeywordsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listKeywords)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")

	return cmd
}

func listKeywords(globalFlags *types.GlobalFlags, flags *listKeywordsFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

