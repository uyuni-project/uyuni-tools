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

type getPillarFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ImageId          int
}

func getPillarCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getPillar",
		Short: "Get pillar data of an image. The "size" entries are converted to string.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getPillarFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getPillar)
		},
	}

	cmd.Flags().String("ImageId", "", "")

	return cmd
}

func getPillar(globalFlags *types.GlobalFlags, flags *getPillarFlags, cmd *cobra.Command, args []string) error {

res, err := image.Image(&flags.ConnectionDetails, flags.ImageId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

