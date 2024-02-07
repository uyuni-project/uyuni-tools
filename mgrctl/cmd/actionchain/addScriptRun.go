package actionchain

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/actionchain"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addScriptRunFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	ChainLabel          string
	ScriptLabel          string
	Uid          string
	Gid          string
	Timeout          int
	ScriptBody          string
}

func addScriptRunCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addScriptRun",
		Short: "Add an action with label to run a script to an Action Chain.
 NOTE: The script body must be Base64 encoded!",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addScriptRunFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addScriptRun)
		},
	}

	cmd.Flags().String("Sid", "", "System ID")
	cmd.Flags().String("ChainLabel", "", "Label of the chain")
	cmd.Flags().String("ScriptLabel", "", "Label of the script")
	cmd.Flags().String("Uid", "", "User ID on the particular system")
	cmd.Flags().String("Gid", "", "Group ID on the particular system")
	cmd.Flags().String("Timeout", "", "Timeout")
	cmd.Flags().String("ScriptBody", "", "Base64 encoded script body")

	return cmd
}

func addScriptRun(globalFlags *types.GlobalFlags, flags *addScriptRunFlags, cmd *cobra.Command, args []string) error {

res, err := actionchain.Actionchain(&flags.ConnectionDetails, flags.Sid, flags.ChainLabel, flags.ScriptLabel, flags.Uid, flags.Gid, flags.Timeout, flags.ScriptBody)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

