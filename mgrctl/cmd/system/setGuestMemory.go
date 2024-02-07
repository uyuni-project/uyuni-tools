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

type setGuestMemoryFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Memory          int
}

func setGuestMemoryCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setGuestMemory",
		Short: "Schedule an action of a guest's host, to set that guest's memory
          allocation",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setGuestMemoryFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setGuestMemory)
		},
	}

	cmd.Flags().String("Sid", "", "The guest's system id")
	cmd.Flags().String("Memory", "", "The amount of memory to          allocate to the guest")

	return cmd
}

func setGuestMemory(globalFlags *types.GlobalFlags, flags *setGuestMemoryFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Memory)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

