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

type listEmptySystemProfilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listEmptySystemProfilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listEmptySystemProfiles",
		Short: "Returns a list of empty system profiles visible to user (created by the createSystemProfile method).",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listEmptySystemProfilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listEmptySystemProfiles)
		},
	}


	return cmd
}

func listEmptySystemProfiles(globalFlags *types.GlobalFlags, flags *listEmptySystemProfilesFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

