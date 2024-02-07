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

type importOSImageFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
	Version          string
	Arch          string
}

func importOSImageCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "importOSImage",
		Short: "Import an image and schedule an inspect afterwards",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags importOSImageFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, importOSImage)
		},
	}

	cmd.Flags().String("Name", "", "image name as specified in the store")
	cmd.Flags().String("Version", "", "version to import")
	cmd.Flags().String("Arch", "", "image architecture")

	return cmd
}

func importOSImage(globalFlags *types.GlobalFlags, flags *importOSImageFlags, cmd *cobra.Command, args []string) error {

res, err := image.Image(&flags.ConnectionDetails, flags.Name, flags.Version, flags.Arch)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

