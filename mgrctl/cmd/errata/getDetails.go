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

type getDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
}

func getDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDetails",
		Short: "Retrieves the details for the erratum matching the given advisory name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDetails)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")

	return cmd
}

func getDetails(globalFlags *types.GlobalFlags, flags *getDetailsFlags, cmd *cobra.Command, args []string) error {

	res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
