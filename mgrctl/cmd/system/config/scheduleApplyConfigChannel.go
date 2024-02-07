package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/config"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type scheduleApplyConfigChannelFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids          []int
	EarliestOccurrence          $date
	Test          bool
}

func scheduleApplyConfigChannelCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleApplyConfigChannel",
		Short: "Schedule highstate application for a given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleApplyConfigChannelFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleApplyConfigChannel)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("EarliestOccurrence", "", "")
	cmd.Flags().String("Test", "", "Run states in test-only mode")

	return cmd
}

func scheduleApplyConfigChannel(globalFlags *types.GlobalFlags, flags *scheduleApplyConfigChannelFlags, cmd *cobra.Command, args []string) error {

res, err := config.Config(&flags.ConnectionDetails, flags.Sids, flags.EarliestOccurrence, flags.Test)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

