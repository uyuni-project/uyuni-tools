package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setGroupMembershipFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
	Sgid                  int
	Member                bool
}

func setGroupMembershipCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setGroupMembership",
		Short: "Set a servers membership in a given group.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setGroupMembershipFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setGroupMembership)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Sgid", "", "")
	cmd.Flags().String("Member", "", "'1' to assign the given server to the given server group, '0' to remove the given server from the given server group.")

	return cmd
}

func setGroupMembership(globalFlags *types.GlobalFlags, flags *setGroupMembershipFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Sgid, flags.Member)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
