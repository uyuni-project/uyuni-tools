package keys

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/keys"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type updateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Description          string
	Type          string
	Content          string
}

func updateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates type and content of the key identified by the description",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, update)
		},
	}

	cmd.Flags().String("Description", "", "")
	cmd.Flags().String("Type", "", "valid values are GPG or SSL")
	cmd.Flags().String("Content", "", "")

	return cmd
}

func update(globalFlags *types.GlobalFlags, flags *updateFlags, cmd *cobra.Command, args []string) error {

res, err := keys.Keys(&flags.ConnectionDetails, flags.Description, flags.Type, flags.Content)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

