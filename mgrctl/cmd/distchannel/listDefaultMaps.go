package distchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/distchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listDefaultMapsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listDefaultMapsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listDefaultMaps",
		Short: "Lists the default distribution channel maps",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listDefaultMapsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listDefaultMaps)
		},
	}


	return cmd
}

func listDefaultMaps(globalFlags *types.GlobalFlags, flags *listDefaultMapsFlags, cmd *cobra.Command, args []string) error {

res, err := distchannel.Distchannel(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

