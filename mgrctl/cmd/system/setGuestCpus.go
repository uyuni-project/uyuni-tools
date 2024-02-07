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

type setGuestCpusFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	NumOfCpus          int
}

func setGuestCpusCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setGuestCpus",
		Short: "Schedule an action of a guest's host, to set that guest's CPU
          allocation",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setGuestCpusFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setGuestCpus)
		},
	}

	cmd.Flags().String("Sid", "", "The guest's system id")
	cmd.Flags().String("NumOfCpus", "", "The number of virtual cpus to          allocate to the guest")

	return cmd
}

func setGuestCpus(globalFlags *types.GlobalFlags, flags *setGuestCpusFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.NumOfCpus)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

