package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type scheduleFileComparisonsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	Path          string
	Sids          []long
}

func scheduleFileComparisonsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleFileComparisons",
		Short: "Schedule a comparison of the latest revision of a file
 against the version deployed on a list of systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleFileComparisonsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleFileComparisons)
		},
	}

	cmd.Flags().String("Label", "", "label of config channel")
	cmd.Flags().String("Path", "", "file path")
	cmd.Flags().String("Sids", "", "$desc")

	return cmd
}

func scheduleFileComparisons(globalFlags *types.GlobalFlags, flags *scheduleFileComparisonsFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.Path, flags.Sids)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

