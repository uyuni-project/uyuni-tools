package delta

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/image/delta"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createDeltaImageFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SourceImageId          int
	TargetImageId          int
	File          string
	Pillar          struct
}

func createDeltaImageCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createDeltaImage",
		Short: "Import an image and schedule an inspect afterwards. The "size" entries in the pillar
 should be passed as string.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createDeltaImageFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createDeltaImage)
		},
	}

	cmd.Flags().String("SourceImageId", "", "")
	cmd.Flags().String("TargetImageId", "", "")
	cmd.Flags().String("File", "", "")
	cmd.Flags().String("Pillar", "", "")

	return cmd
}

func createDeltaImage(globalFlags *types.GlobalFlags, flags *createDeltaImageFlags, cmd *cobra.Command, args []string) error {

res, err := delta.Delta(&flags.ConnectionDetails, flags.SourceImageId, flags.TargetImageId, flags.File, flags.Pillar)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

