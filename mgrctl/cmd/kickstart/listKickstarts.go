package kickstart

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listKickstartsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listKickstartsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listKickstarts",
		Short: "Provides a list of kickstart profiles visible to the user's
 org",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listKickstartsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listKickstarts)
		},
	}


	return cmd
}

func listKickstarts(globalFlags *types.GlobalFlags, flags *listKickstartsFlags, cmd *cobra.Command, args []string) error {

res, err := kickstart.Kickstart(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

