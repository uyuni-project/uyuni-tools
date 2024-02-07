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

type scheduleImageBuildFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProfileLabel          string
	Version          string
	BuildHostId          int
	EarliestOccurrence          $date
}

func scheduleImageBuildCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleImageBuild",
		Short: "Schedule an image build",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleImageBuildFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleImageBuild)
		},
	}

	cmd.Flags().String("ProfileLabel", "", "")
	cmd.Flags().String("Version", "", "version to build or empty")
	cmd.Flags().String("BuildHostId", "", "system id of the build host")
	cmd.Flags().String("EarliestOccurrence", "", "earliest the build can run.")

	return cmd
}

func scheduleImageBuild(globalFlags *types.GlobalFlags, flags *scheduleImageBuildFlags, cmd *cobra.Command, args []string) error {

res, err := image.Image(&flags.ConnectionDetails, flags.ProfileLabel, flags.Version, flags.BuildHostId, flags.EarliestOccurrence)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

