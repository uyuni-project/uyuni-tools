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

type addImageFileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ImageId          int
	File          string
	Type          string
	External          bool
}

func addImageFileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addImageFile",
		Short: "Delete image file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addImageFileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addImageFile)
		},
	}

	cmd.Flags().String("ImageId", "", "ID of the image")
	cmd.Flags().String("File", "", "the file name, it must exist in the store")
	cmd.Flags().String("Type", "", "the image type")
	cmd.Flags().String("External", "", "the file is external")

	return cmd
}

func addImageFile(globalFlags *types.GlobalFlags, flags *addImageFileFlags, cmd *cobra.Command, args []string) error {

res, err := image.Image(&flags.ConnectionDetails, flags.ImageId, flags.File, flags.Type, flags.External)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

