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

type listPackagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ImageId               int
}

func listPackagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listPackages",
		Short: "List the installed packages on the given image",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listPackagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listPackages)
		},
	}

	cmd.Flags().String("ImageId", "", "")

	return cmd
}

func listPackages(globalFlags *types.GlobalFlags, flags *listPackagesFlags, cmd *cobra.Command, args []string) error {

	res, err := image.Image(&flags.ConnectionDetails, flags.ImageId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
