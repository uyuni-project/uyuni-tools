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

type getChannelLastBuildByIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Id          int
}

func getChannelLastBuildByIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getChannelLastBuildById",
		Short: "Returns the last build date of the repomd.xml file
 for the given channel as a localised string.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getChannelLastBuildByIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getChannelLastBuildById)
		},
	}

	cmd.Flags().String("Id", "", "id of channel wanted")

	return cmd
}

func getChannelLastBuildById(globalFlags *types.GlobalFlags, flags *getChannelLastBuildByIdFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Id)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

