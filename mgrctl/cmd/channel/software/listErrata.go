package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listErrataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	StartDate          $type
	EndDate          $type
	LastModified          bool
}

func listErrataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listErrata",
		Short: "List the errata applicable to a channel after given startDate",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listErrataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listErrata)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "channel to query")
	cmd.Flags().String("StartDate", "", "")
	cmd.Flags().String("EndDate", "", "")
	cmd.Flags().String("LastModified", "", "select by last modified or not")

	return cmd
}

func listErrata(globalFlags *types.GlobalFlags, flags *listErrataFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.StartDate, flags.EndDate, flags.LastModified)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

