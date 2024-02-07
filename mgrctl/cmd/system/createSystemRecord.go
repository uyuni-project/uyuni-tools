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

type createSystemRecordFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	KsLabel          string
	SystemName          string
	KOptions          string
	Comment          string
	$param.getFlagName()          $param.getType()
}

func createSystemRecordCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createSystemRecord",
		Short: "Creates a cobbler system record with the specified kickstart label",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createSystemRecordFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createSystemRecord)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("KsLabel", "", "")
	cmd.Flags().String("SystemName", "", "")
	cmd.Flags().String("KOptions", "", "")
	cmd.Flags().String("Comment", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func createSystemRecord(globalFlags *types.GlobalFlags, flags *createSystemRecordFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.KsLabel, flags.SystemName, flags.KOptions, flags.Comment, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

