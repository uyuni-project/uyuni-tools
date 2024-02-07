package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/image/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getCustomValuesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
}

func getCustomValuesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getCustomValues",
		Short: "Get the custom data values defined for the image profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getCustomValuesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getCustomValues)
		},
	}

	cmd.Flags().String("Label", "", "")

	return cmd
}

func getCustomValues(globalFlags *types.GlobalFlags, flags *getCustomValuesFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
