package virtualhostmanager

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/virtualhostmanager"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getDetailFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
}

func getDetailCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDetail",
		Short: "Gets details of a Virtual Host Manager with a given label",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDetailFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDetail)
		},
	}

	cmd.Flags().String("Label", "", "Virtual Host Manager label")

	return cmd
}

func getDetail(globalFlags *types.GlobalFlags, flags *getDetailFlags, cmd *cobra.Command, args []string) error {

res, err := virtualhostmanager.Virtualhostmanager(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

