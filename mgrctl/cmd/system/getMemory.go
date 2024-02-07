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

type getMemoryFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func getMemoryCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getMemory",
		Short: "Gets the memory information for a system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getMemoryFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getMemory)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getMemory(globalFlags *types.GlobalFlags, flags *getMemoryFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
