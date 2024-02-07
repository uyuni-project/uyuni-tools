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

type setVariablesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Netboot          bool
	$param.getFlagName()          $param.getType()
}

func setVariablesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setVariables",
		Short: "Sets a list of kickstart variables in the cobbler system record
 for the specified server.
  Note: This call assumes that a system record exists in cobbler for the
  given system and will raise an XMLRPC fault if that is not the case.
  To create a system record over xmlrpc use system.createSystemRecord

  To create a system record in the Web UI  please go to
  System -&gt; &lt;Specified System&gt; -&gt; Provisioning -&gt;
  Select a Kickstart profile -&gt; Create Cobbler System Record.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setVariablesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setVariables)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Netboot", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func setVariables(globalFlags *types.GlobalFlags, flags *setVariablesFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Netboot, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

