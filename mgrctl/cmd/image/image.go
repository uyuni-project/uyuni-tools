package image

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Provides methods to access and modify images.",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(listImagesCommand(globalFlags))
	cmd.AddCommand(getPillarCommand(globalFlags))
	cmd.AddCommand(listPackagesCommand(globalFlags))
	cmd.AddCommand(scheduleImageBuildCommand(globalFlags))
	cmd.AddCommand(getRelevantErrataCommand(globalFlags))
	cmd.AddCommand(getCustomValuesCommand(globalFlags))
	cmd.AddCommand(setPillarCommand(globalFlags))
	cmd.AddCommand(importContainerImageCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(importOSImageCommand(globalFlags))
	cmd.AddCommand(deleteImageFileCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(addImageFileCommand(globalFlags))

	return cmd
}
