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

type getRelevantErrataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ImageId               int
}

func getRelevantErrataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getRelevantErrata",
		Short: "Returns a list of all errata that are relevant for the image",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getRelevantErrataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getRelevantErrata)
		},
	}

	cmd.Flags().String("ImageId", "", "")

	return cmd
}

func getRelevantErrata(globalFlags *types.GlobalFlags, flags *getRelevantErrataFlags, cmd *cobra.Command, args []string) error {

	res, err := image.Image(&flags.ConnectionDetails, flags.ImageId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
