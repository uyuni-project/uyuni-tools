package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/image/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listImageProfilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listImageProfilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listImageProfiles",
		Short: "List available image profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listImageProfilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listImageProfiles)
		},
	}

	return cmd
}

func listImageProfiles(globalFlags *types.GlobalFlags, flags *listImageProfilesFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
