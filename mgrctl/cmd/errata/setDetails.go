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

type setDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
}

func setDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setDetails",
		Short: "Set erratum details. All arguments are optional and will only be modified
 if included in the struct. This method will only allow for modification of custom
 errata created either through the UI or API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setDetails)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")

	return cmd
}

func setDetails(globalFlags *types.GlobalFlags, flags *setDetailsFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

