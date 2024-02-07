package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
	Id                    int
}

func getDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDetails",
		Short: "Lookup config channel details.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDetails)
		},
	}

	cmd.Flags().String("Label", "", "")
	cmd.Flags().String("Id", "", "the channel ID")

	return cmd
}

func getDetails(globalFlags *types.GlobalFlags, flags *getDetailsFlags, cmd *cobra.Command, args []string) error {

	res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.Id)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
