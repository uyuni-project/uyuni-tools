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

type listMapsForOrgFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId                 int
}

func listMapsForOrgCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listMapsForOrg",
		Short: "Lists distribution channel maps valid for the user's organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listMapsForOrgFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listMapsForOrg)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func listMapsForOrg(globalFlags *types.GlobalFlags, flags *listMapsForOrgFlags, cmd *cobra.Command, args []string) error {

	res, err := distchannel.Distchannel(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
