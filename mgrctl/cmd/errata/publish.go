package errata

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/errata"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type publishFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
	$param.getFlagName()          $param.getType()
}

func publishCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Adds an existing errata to a set of channels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags publishFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, publish)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func publish(globalFlags *types.GlobalFlags, flags *publishFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

