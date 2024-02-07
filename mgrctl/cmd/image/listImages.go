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

type listImagesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listImagesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listImages",
		Short: "List available images",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listImagesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listImages)
		},
	}


	return cmd
}

func listImages(globalFlags *types.GlobalFlags, flags *listImagesFlags, cmd *cobra.Command, args []string) error {

res, err := image.Image(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

