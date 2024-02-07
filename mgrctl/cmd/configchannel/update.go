package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type updateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	Name          string
	Description          string
}

func updateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a global config channel. Caller must be at least a
 config admin or an organization admin, or have access to a system containing this
 config channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, update)
		},
	}

	cmd.Flags().String("Label", "", "")
	cmd.Flags().String("Name", "", "")
	cmd.Flags().String("Description", "", "")

	return cmd
}

func update(globalFlags *types.GlobalFlags, flags *updateFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.Name, flags.Description)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

