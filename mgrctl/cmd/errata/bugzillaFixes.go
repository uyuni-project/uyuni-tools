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

type bugzillaFixesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
}

func bugzillaFixesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bugzillaFixes",
		Short: "Get the Bugzilla fixes for an erratum matching the given
 advisoryName. The bugs will be returned in a struct where the bug id is
 the key.  i.e. 208144="errata.bugzillaFixes Method Returns different
 results than docs say"
 For those errata that are present in both vendor and user organizations under the same advisory name,
 this method retrieves the list of Bugzilla fixes of both of them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags bugzillaFixesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, bugzillaFixes)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")

	return cmd
}

func bugzillaFixes(globalFlags *types.GlobalFlags, flags *bugzillaFixesFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

