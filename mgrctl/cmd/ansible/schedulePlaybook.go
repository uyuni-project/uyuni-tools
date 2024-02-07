package ansible

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/ansible"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type schedulePlaybookFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	PlaybookPath          string
	InventoryPath          string
	ControlNodeId          int
	EarliestOccurrence          $date
	ActionChainLabel          string
	TestMode          bool
	$param.getFlagName()          $param.getType()
}

func schedulePlaybookCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedulePlaybook",
		Short: "Schedule a playbook execution",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags schedulePlaybookFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, schedulePlaybook)
		},
	}

	cmd.Flags().String("PlaybookPath", "", "")
	cmd.Flags().String("InventoryPath", "", "path to Ansible inventory or empty")
	cmd.Flags().String("ControlNodeId", "", "system ID of the control node")
	cmd.Flags().String("EarliestOccurrence", "", "earliest the execution command can be sent to the control node. ignored when actionChainLabel is used")
	cmd.Flags().String("ActionChainLabel", "", "label of an action chain to use, or None")
	cmd.Flags().String("TestMode", "", "'true' if the playbook shall be executed in test mode")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func schedulePlaybook(globalFlags *types.GlobalFlags, flags *schedulePlaybookFlags, cmd *cobra.Command, args []string) error {

res, err := ansible.Ansible(&flags.ConnectionDetails, flags.PlaybookPath, flags.InventoryPath, flags.ControlNodeId, flags.EarliestOccurrence, flags.ActionChainLabel, flags.TestMode, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

