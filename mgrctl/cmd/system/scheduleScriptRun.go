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

type scheduleScriptRunFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	$param.getFlagName()          $param.getType()
	Username          string
	Groupname          string
	Timeout          int
	Script          string
	EarliestOccurrence          $date
	Sid          int
}

func scheduleScriptRunCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleScriptRun",
		Short: "Schedule a script to run.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleScriptRunFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleScriptRun)
		},
	}

	cmd.Flags().String("Label", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("Username", "", "User to run script as.")
	cmd.Flags().String("Groupname", "", "Group to run script as.")
	cmd.Flags().String("Timeout", "", "Seconds to allow the script to runbefore timing out.")
	cmd.Flags().String("Script", "", "Contents of the script to run. Must start with a shebang (e.g. #!/bin/bash)")
	cmd.Flags().String("EarliestOccurrence", "", "Earliest the script can run.")
	cmd.Flags().String("Sid", "", "ID of the server to run the script on.")

	return cmd
}

func scheduleScriptRun(globalFlags *types.GlobalFlags, flags *scheduleScriptRunFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Label, flags.$param.getFlagName(), flags.Username, flags.Groupname, flags.Timeout, flags.Script, flags.EarliestOccurrence, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

