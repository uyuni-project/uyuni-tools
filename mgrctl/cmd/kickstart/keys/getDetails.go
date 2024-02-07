package keys

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/keys"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Description           string
}

func getDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDetails",
		Short: "returns all the data associated with the given key",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDetails)
		},
	}

	cmd.Flags().String("Description", "", "")

	return cmd
}

func getDetails(globalFlags *types.GlobalFlags, flags *getDetailsFlags, cmd *cobra.Command, args []string) error {

	res, err := keys.Keys(&flags.ConnectionDetails, flags.Description)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
