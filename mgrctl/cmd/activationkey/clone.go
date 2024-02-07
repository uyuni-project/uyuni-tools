package activationkey

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/activationkey"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type cloneFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Key                   string
	CloneDescription      string
}

func cloneCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone an existing activation key.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags cloneFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, clone)
		},
	}

	cmd.Flags().String("Key", "", "Key to be cloned.")
	cmd.Flags().String("CloneDescription", "", "Description of the cloned key.")

	return cmd
}

func clone(globalFlags *types.GlobalFlags, flags *cloneFlags, cmd *cobra.Command, args []string) error {

	res, err := activationkey.Activationkey(&flags.ConnectionDetails, flags.Key, flags.CloneDescription)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
