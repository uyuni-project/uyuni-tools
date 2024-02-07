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

type findByCveFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	CveName          string
}

func findByCveCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "findByCve",
		Short: "Lookup the details for errata associated with the given CVE
 (e.g. CVE-2008-3270)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags findByCveFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, findByCve)
		},
	}

	cmd.Flags().String("CveName", "", "")

	return cmd
}

func findByCve(globalFlags *types.GlobalFlags, flags *findByCveFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.CveName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

