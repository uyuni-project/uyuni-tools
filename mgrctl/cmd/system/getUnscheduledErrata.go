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

type getUnscheduledErrataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func getUnscheduledErrataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getUnscheduledErrata",
		Short: "Provides an array of errata that are applicable to a given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getUnscheduledErrataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getUnscheduledErrata)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getUnscheduledErrata(globalFlags *types.GlobalFlags, flags *getUnscheduledErrataFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

