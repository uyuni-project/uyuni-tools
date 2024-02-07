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

type listCvesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
}

func listCvesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listCves",
		Short: "Returns a list of http://cve.mitre.org/_blankCVEs applicable to the errata
 with the given advisory name. For those errata that are present in both vendor and user organizations under the
 same advisory name, this method retrieves the list of CVEs of both of them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listCvesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listCves)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")

	return cmd
}

func listCves(globalFlags *types.GlobalFlags, flags *listCvesFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

