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

type publishAsOriginalFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
	$param.getFlagName()          $param.getType()
}

func publishAsOriginalCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publishAsOriginal",
		Short: "Adds an existing cloned errata to a set of cloned
 channels according to its original erratum",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags publishAsOriginalFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, publishAsOriginal)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func publishAsOriginal(globalFlags *types.GlobalFlags, flags *publishAsOriginalFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

