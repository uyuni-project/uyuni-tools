package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type removeScriptFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	ScriptId          int
}

func removeScriptCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeScript",
		Short: "Remove a script from a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeScriptFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeScript)
		},
	}

	cmd.Flags().String("KsLabel", "", "The kickstart from which to remove the script from.")
	cmd.Flags().String("ScriptId", "", "The id of the script to remove.")

	return cmd
}

func removeScript(globalFlags *types.GlobalFlags, flags *removeScriptFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.ScriptId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

