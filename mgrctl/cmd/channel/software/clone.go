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

type cloneFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OriginalLabel          string
	OriginalState          bool
}

func cloneCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone a channel.  If arch_label is omitted, the arch label of the
      original channel will be used. If parent_label is omitted, the clone will be
      a base channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags cloneFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, clone)
		},
	}

	cmd.Flags().String("OriginalLabel", "", "")
	cmd.Flags().String("OriginalState", "", "")

	return cmd
}

func clone(globalFlags *types.GlobalFlags, flags *cloneFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.OriginalLabel, flags.OriginalState)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

