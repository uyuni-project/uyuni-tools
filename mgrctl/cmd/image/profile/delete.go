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

type deleteFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
}

func deleteCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an image profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, delete)
		},
	}

	cmd.Flags().String("Label", "", "")

	return cmd
}

func delete(globalFlags *types.GlobalFlags, flags *deleteFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
