package snippet

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/snippet"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listDefaultFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listDefaultCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listDefault",
		Short: "List only pre-made default snippets for the logged in user.
    These snipppets are not editable.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listDefaultFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listDefault)
		},
	}


	return cmd
}

func listDefault(globalFlags *types.GlobalFlags, flags *listDefaultFlags, cmd *cobra.Command, args []string) error {

res, err := snippet.Snippet(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

