package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getSoftwareDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getSoftwareDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getSoftwareDetails",
		Short: "Gets kickstart profile software details.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getSoftwareDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getSoftwareDetails)
		},
	}

	cmd.Flags().String("KsLabel", "", "the label of the kickstart profile")

	return cmd
}

func getSoftwareDetails(globalFlags *types.GlobalFlags, flags *getSoftwareDetailsFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

