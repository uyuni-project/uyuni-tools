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

type setPillarFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ImageId          int
	PillarData          struct
}

func setPillarCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setPillar",
		Short: "Set pillar data of an image. The "size" entries should be passed as string.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setPillarFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setPillar)
		},
	}

	cmd.Flags().String("ImageId", "", "")
	cmd.Flags().String("PillarData", "", "")

	return cmd
}

func setPillar(globalFlags *types.GlobalFlags, flags *setPillarFlags, cmd *cobra.Command, args []string) error {

res, err := image.Image(&flags.ConnectionDetails, flags.ImageId, flags.PillarData)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

