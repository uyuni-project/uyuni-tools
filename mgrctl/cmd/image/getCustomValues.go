package image

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/image"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getCustomValuesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ImageId               int
}

func getCustomValuesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getCustomValues",
		Short: "Get the custom data values defined for the image",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getCustomValuesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getCustomValues)
		},
	}

	cmd.Flags().String("ImageId", "", "")

	return cmd
}

func getCustomValues(globalFlags *types.GlobalFlags, flags *getCustomValuesFlags, cmd *cobra.Command, args []string) error {

	res, err := image.Image(&flags.ConnectionDetails, flags.ImageId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
