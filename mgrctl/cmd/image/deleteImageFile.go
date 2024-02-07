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

type deleteImageFileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ImageId               int
	File                  string
}

func deleteImageFileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteImageFile",
		Short: "Delete image file",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteImageFileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteImageFile)
		},
	}

	cmd.Flags().String("ImageId", "", "ID of the image")
	cmd.Flags().String("File", "", "the file name")

	return cmd
}

func deleteImageFile(globalFlags *types.GlobalFlags, flags *deleteImageFileFlags, cmd *cobra.Command, args []string) error {

	res, err := image.Image(&flags.ConnectionDetails, flags.ImageId, flags.File)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
