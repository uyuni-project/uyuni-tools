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

type createOrUpdateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
	Contents          string
}

func createOrUpdateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createOrUpdate",
		Short: "Will create a snippet with the given name and contents if it
      doesn't exist. If it does exist, the existing snippet will be updated.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createOrUpdateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createOrUpdate)
		},
	}

	cmd.Flags().String("Name", "", "")
	cmd.Flags().String("Contents", "", "")

	return cmd
}

func createOrUpdate(globalFlags *types.GlobalFlags, flags *createOrUpdateFlags, cmd *cobra.Command, args []string) error {

res, err := snippet.Snippet(&flags.ConnectionDetails, flags.Name, flags.Contents)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

