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

type importContainerImageFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
	Version          string
	BuildHostId          int
	StoreLabel          string
	ActivationKey          string
	EarliestOccurrence          $date
}

func importContainerImageCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "importContainerImage",
		Short: "Import an image and schedule an inspect afterwards",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags importContainerImageFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, importContainerImage)
		},
	}

	cmd.Flags().String("Name", "", "image name as specified in the store")
	cmd.Flags().String("Version", "", "version to import or empty")
	cmd.Flags().String("BuildHostId", "", "system ID of the build host")
	cmd.Flags().String("StoreLabel", "", "")
	cmd.Flags().String("ActivationKey", "", "activation key to get the channel data from")
	cmd.Flags().String("EarliestOccurrence", "", "earliest the following inspect can run")

	return cmd
}

func importContainerImage(globalFlags *types.GlobalFlags, flags *importContainerImageFlags, cmd *cobra.Command, args []string) error {

res, err := image.Image(&flags.ConnectionDetails, flags.Name, flags.Version, flags.BuildHostId, flags.StoreLabel, flags.ActivationKey, flags.EarliestOccurrence)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

