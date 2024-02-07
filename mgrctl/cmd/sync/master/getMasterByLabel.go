package master

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/master"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getMasterByLabelFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label                 string
}

func getMasterByLabelCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getMasterByLabel",
		Short: "Find a Master by specifying its label",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getMasterByLabelFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getMasterByLabel)
		},
	}

	cmd.Flags().String("Label", "", "Label of the desired Master")

	return cmd
}

func getMasterByLabel(globalFlags *types.GlobalFlags, flags *getMasterByLabelFlags, cmd *cobra.Command, args []string) error {

	res, err := master.Master(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
