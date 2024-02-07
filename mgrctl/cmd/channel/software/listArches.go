package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listArchesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listArchesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listArches",
		Short: "Lists the potential software channel architectures that can be created",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listArchesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listArches)
		},
	}

	return cmd
}

func listArches(globalFlags *types.GlobalFlags, flags *listArchesFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
